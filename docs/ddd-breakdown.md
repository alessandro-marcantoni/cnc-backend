# Domain-Driven Design for Circolo Nautico Cattolica

## 1. Ubiquitous Language Analysis

Key concepts emerging from the document:

- **Socio (Member)**: club member with membership card
- **Tessera Socio (Membership Card)**: represents active membership
- **Servizio (Service)**: rentable resource from the club
- **Affitto (Rental)**: relationship between member and service with payment status
- **Lista di Attesa (Waiting List)**: queue for unavailable services
- **Diritto di Prelazione (Right of Preemption)**: member's priority for renewal
- **Esclusione (Exclusion)**: definitive removal of member

## 2. Bounded Contexts

### **Membership Context**
Manages the lifecycle of members and membership cards.

### **Service Rental Context** (Core Domain)
Manages service rentals, availability, and waiting lists.

### **Payment Context**
Manages payment status for membership cards and services.

### **Course Management Context** (Future)
For managing courses and students.

## 3. Tactical Design - Service Rental Context

This is the core domain. Here's the detailed design:

### **Aggregates**

#### **Aggregate: Member**
```
Member (Aggregate Root)
├─ MemberId (Value Object)
├─ PersonalInfo (Value Object)
│  ├─ name
│  ├─ surname
│  ├─ email
│  └─ phone
├─ MembershipStatus (Value Object)
│  ├─ year
│  ├─ status (Active, Expired, ExclusionDeliberated)
│  ├─ expirationDate
│  └─ isPaid (boolean)
└─ Methods:
   ├─ renewMembership(year)
   ├─ deliberateExclusion(reason)
   ├─ removeDefinitively()
   └─ canRentServices() -> boolean

Invariants:
- A member can be removed only if status is ExclusionDeliberated
- A member with expired membership retains status until deliberation
- Even non-paying members remain members until exclusion
```

#### **Aggregate: RentableService**
```
RentableService (Aggregate Root)
├─ ServiceId (Value Object)
├─ ServiceType (Value Object/Enum)
│  ├─ category (Numbered/ByBoatType)
│  ├─ name
│  ├─ hasDiscount (if box rented)
│  └─ pricingType (Fixed/ByBoatSize)
├─ ServiceNumber (Value Object) - optional for numbered services
├─ AvailabilityManager
│  ├─ totalCapacity (for numbered services)
│  ├─ currentRentals (Set<RentalId>)
│  └─ calculateAvailable() -> int
├─ Price (Value Object)
│  ├─ baseAmount
│  └─ boatSizeVariants (Map<BoatSize, Amount>) - optional
└─ Methods:
   ├─ isAvailable() -> boolean
   ├─ calculatePriceFor(member, boatSize?) -> Money
   ├─ applyBoxDiscount(price) -> Money
   └─ reserveFor(memberId) -> Result

Service Types (Enum):
- OPEN_BOARD_RACK (Rastrelliera tavole aperta)
- DINGHY_BOAT_SPACES (Posti barca piazzale derive)
- BOX
- CLOSED_BOARD_RACK (Rastrelliera tavole chiusa) - with box discount
- LARGE_WOODEN_LOCKERS (Armadietti grandi legno)
- NORMAL_LOCKERS (Armadietti normali)
- SURF_STORAGE_WORKSHOP (Deposito surf lato officina) - with box discount
- CLOSED_SUP_STORAGE (Deposito SUP chiuso)
- EXTERNAL_CANOE_RACK (Rastrelliera canoe esterna) - with box discount
- VENTENA_RIVER_SPACES (Posti a terra fiume Ventena) - by boat type
- TAVOLLO_RIVER_SPACES (Posti barca fiume Tavollo) - by boat type

Invariants:
- Numbered services have fixed capacity
- River spaces availability varies by boat size/type
- Discount services require an active box rental
```

#### **Aggregate: ServiceRental**
```
ServiceRental (Aggregate Root)
├─ RentalId (Value Object)
├─ MemberId (Value Object)
├─ ServiceId (Value Object)
├─ ServiceNumber (Value Object) - optional
├─ RentalPeriod (Value Object)
│  ├─ year
│  ├─ startDate
│  └─ endDate
├─ PaymentStatus (Value Object)
│  ├─ status (Paid/Unpaid/Late)
│  ├─ amount
│  ├─ dueDate
│  └─ paidDate (optional)
├─ HasPreemptionRight (boolean)
├─ BoatInfo (Value Object) - optional for boat-based services
│  ├─ boatType
│  └─ size
└─ Methods:
   ├─ markAsPaid(date, amount)
   ├─ checkIfLate(currentDate) -> boolean
   ├─ renew(newYear) -> ServiceRental
   ├─ terminate()
   └─ grantPreemptionRight()

Invariants:
- A rental must have a valid member and service
- Payment status must be updated only with valid dates
- Preemption right is automatic for existing rentals
- Renewal creates a new rental for the next year
```

#### **Aggregate: WaitingList**
```
WaitingList (Aggregate Root)
├─ WaitingListId (Value Object)
├─ ServiceId (Value Object)
├─ Entries (List<WaitingListEntry>)
│  └─ WaitingListEntry (Entity)
│     ├─ entryId
│     ├─ memberId
│     ├─ requestDate
│     ├─ position
│     ├─ preferredNumber (optional)
│     └─ boatInfo (optional)
└─ Methods:
   ├─ addMember(memberId, preferences)
   ├─ removeMember(memberId)
   ├─ getNextInLine() -> WaitingListEntry
   ├─ notifyAvailability(serviceNumber?)
   └─ reorderByDate()

Invariants:
- Entries are ordered by request date (FIFO)
- A member can appear only once per service waiting list
- When service becomes available, first in line is notified
```

### **Domain Services**

#### **RentalManagementService**
```
RentalManagementService
└─ Methods:
   ├─ rentService(memberId, serviceId, boatInfo?, preferredNumber?)
   │  -> Result<ServiceRental, Error>
   │
   ├─ addToWaitingList(memberId, serviceId, preferences)
   │  -> Result<WaitingListEntry, Error>
   │
   ├─ processRenewal(rentalId)
   │  -> Result<ServiceRental, Error>
   │
   ├─ releaseService(rentalId)
   │  -> notifies waiting list
   │
   └─ calculateDiscountedPrice(memberId, serviceId)
      -> Money

Business Rules:
- Check if member has active membership
- Check if service is available
- Apply box discount if applicable
- Verify preemption rights during renewal
- Process waiting list when service becomes available
```

#### **AvailabilityCalculationService**
```
AvailabilityCalculationService
└─ Methods:
   ├─ getAvailableServices()
   │  -> List<ServiceAvailability>
   │
   ├─ getAvailableNumbersFor(serviceId)
   │  -> List<ServiceNumber>
   │
   ├─ calculateRiverSpaceAvailability(serviceId, boatSize)
   │  -> int (empirical calculation as mentioned in Q3)
   │
   └─ canAccommodateBoat(serviceId, boatInfo)
      -> boolean

Note: For river spaces, availability is determined empirically
based on renewals and requests (as per Q3)
```

#### **PaymentTrackingService**
```
PaymentTrackingService
└─ Methods:
   ├─ getUnpaidRentals(asOfDate)
   │  -> List<ServiceRental>
   │
   ├─ getLatePayments(asOfDate)
   │  -> List<ServiceRental>
   │
   ├─ getUnpaidMemberships(year)
   │  -> List<Member>
   │
   └─ recordPayment(rentalId, amount, date)
      -> Result<void, Error>
```

### **Domain Events**

```
Events:
├─ MemberRegistered
├─ MembershipRenewed
├─ MemberExclusionDeliberated
├─ MemberRemovedDefinitively
├─ ServiceRented
├─ ServiceReleased
├─ RentalPaymentReceived
├─ RentalPaymentOverdue
├─ MemberAddedToWaitingList
├─ ServiceBecameAvailable
└─ RentalRenewed

These events enable:
- Asynchronous processing
- Notification systems
- Audit trails
- Integration with other contexts
```

### **Value Objects**

```
Key Value Objects:

Money
├─ amount (BigDecimal)
└─ currency

BoatSize (Enum)
├─ SMALL
├─ MEDIUM
└─ LARGE

ServiceNumber
├─ number (int)
└─ isAssigned (boolean)

RentalPeriod
├─ year
├─ startDate
└─ endDate
└─ Methods:
   ├─ isActive(date) -> boolean
   ├─ isExpired(date) -> boolean
   └─ daysRemaining(date) -> int

PaymentStatus
├─ status (Enum: Paid, Unpaid, Late)
├─ amount
├─ dueDate
└─ paidDate

PersonalInfo
├─ name
├─ surname
├─ email
└─ phone
```

## 4. Repository Interfaces

```
Repositories (to be implemented by infrastructure):

MemberRepository
├─ findById(memberId) -> Member
├─ findAll() -> List<Member>
├─ findByMembershipStatus(status) -> List<Member>
├─ findUnpaidMemberships(year) -> List<Member>
├─ save(member)
└─ delete(memberId)

RentableServiceRepository
├─ findById(serviceId) -> RentableService
├─ findByType(serviceType) -> List<RentableService>
├─ findAvailable() -> List<RentableService>
└─ save(service)

ServiceRentalRepository
├─ findById(rentalId) -> ServiceRental
├─ findByMember(memberId) -> List<ServiceRental>
├─ findByService(serviceId) -> List<ServiceRental>
├─ findUnpaid() -> List<ServiceRental>
├─ findLate(asOfDate) -> List<ServiceRental>
├─ findActiveRentals(year) -> List<ServiceRental>
├─ save(rental)
└─ delete(rentalId)

WaitingListRepository
├─ findByService(serviceId) -> WaitingList
├─ findByMember(memberId) -> List<WaitingListEntry>
└─ save(waitingList)
```

## 5. Application Services (Use Cases)

```
Member Management:
├─ RegisterNewMember
├─ RenewMembershipCard
├─ DeliberateMemberExclusion
├─ RemoveMemberDefinitively
└─ GetMemberRentals

Service Rental Management:
├─ RentServiceToMember
├─ ReleaseServiceRental
├─ RenewServiceRental
├─ RecordPayment
├─ UpdateServiceNumber
└─ AddMemberToWaitingList

Reporting:
├─ GetAvailableServices
├─ GetMemberList
├─ GetUnpaidMemberships
├─ GetUnpaidRentals
├─ GetLatePayments
└─ GetWaitingLists
```

## 6. Key Business Rules Implementation

### **Preemption Right Rule**
```
When processing renewals:
1. Existing rental holders have priority
2. They can keep the same service number
3. They must be processed before new requests
4. Implementation: Check rental history before assigning services
```

### **Box Discount Rule**
```
Services eligible for discount:
- Closed Board Rack
- Surf Storage Workshop
- External Canoe Rack

Calculation:
IF member has active box rental
  THEN apply discount to eligible services
ELSE charge full price
```

### **Member Status Rule**
```
Member lifecycle:
Active → Unpaid → Exclusion Deliberated → Removed

Transitions:
- Active: can rent services
- Unpaid: remains member, keeps rentals, but flagged
- Exclusion Deliberated: by board decision
- Removed: only after deliberation
```

### **Availability Calculation**
```
For numbered services:
  available = totalCapacity - activeRentals

For river spaces (empirical):
  available = estimatedCapacity(boatSizes) - activeRentals
  (manual adjustment based on actual space used)
```

## 7. Context Map

```
┌─────────────────────────┐
│  Membership Context     │
│  (Upstream)             │
│  - Member lifecycle     │
│  - Membership cards     │
└───────────┬─────────────┘
            │ Published Events
            │ (MemberRegistered, etc.)
            ↓
┌─────────────────────────┐
│ Service Rental Context  │◄─────── Anticorruption Layer
│ (Core Domain)           │
│ - Rentals               │
│ - Services              │
│ - Waiting lists         │
└───────────┬─────────────┘
            │ Payment requests
            ↓
┌─────────────────────────┐
│  Payment Context        │
│  (Downstream)           │
│  - Payment processing   │
│  - Payment status       │
└─────────────────────────┘

Future:
┌─────────────────────────┐
│ Course Management       │
│ - Course registration   │
│ - Course payments       │
│ - Instructor management │
└─────────────────────────┘
```

## 8. Strategic Design Considerations

### **Core Domain**: Service Rental Context
This is where the club's competitive advantage and complexity lie. Focus development effort here.

### **Supporting Subdomains**:
- Membership Context
- Payment Context

### **Generic Subdomains**:
- Course Management (can use off-the-shelf solutions)

### **Partnership Patterns**:
- **Shared Kernel**: Between Service Rental and Payment contexts for payment status
- **Customer-Supplier**: Membership Context (supplier) → Service Rental Context (customer)
- **Conformist**: Service Rental Context consumes Member data as-is

## 9. Implementation Recommendations

### **Phase 1**: Core functionality
- Member aggregate with basic CRUD
- Service aggregate with availability
- Rental aggregate with payment tracking
- Basic repositories

### **Phase 2**: Advanced features
- Waiting lists
- Preemption rights
- Discount calculations
- Reporting queries

### **Phase 3**: Future extensions
- Course management integration
- Payment gateway integration
- Automated notifications

### **Technical Considerations**:
- Use **Event Sourcing** for rentals to maintain history (addresses Q6 about history)
- Implement **CQRS** for complex queries (reports, availability)
- Use **Domain Events** for cross-context communication
- Consider **manual availability override** for river spaces (Q3.1)

## 10. Addressing Specific Questions from Document

**Q3/Q3.1**: Availability calculation
- For numbered services: automatic calculation
- For river spaces: empirical calculation with manual override capability
- System should allow administrators to manually adjust available spaces

**Q4**: Preemption rights apply to ALL services, not just boat spaces
- Implemented through `HasPreemptionRight` flag in ServiceRental
- Renewal process checks this flag first

**Q6**: History management
- Use Event Sourcing for ServiceRental aggregate
- Maintain history primarily for continuity (same number/space)
- Full audit trail available through event stream

---

## Conclusion

This DDD design provides a solid foundation that respects domain complexity while remaining flexible for future extensions like course management. The design emphasizes:

- Clear bounded contexts with well-defined responsibilities
- Rich domain models with business logic encapsulated in aggregates
- Domain events for decoupling and audit trails
- Flexibility for empirical calculations (river spaces)
- Support for complex business rules (preemption, discounts)
- Extensibility for future requirements (courses)

The tactical patterns (aggregates, value objects, domain services) ensure that business rules are expressed clearly in code, while the strategic patterns (bounded contexts, context map) provide a scalable architecture for growth.
