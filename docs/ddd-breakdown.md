# Circolo Nautico Cattolica - Domain-Driven Design Model

## Bounded Contexts

### 1. Membership Context
Manages members, membership cards and their status.

### 2. Services Context
Manages offered services, their availability and rentals.

### 3. Courses Context (future extension)
Manages private courses and summer camp courses.

---

## Ubiquitous Language

- **Member (Socio)**: Person with valid membership card or not yet excluded by the board
- **Membership Card (Tessera Socio)**: Annual subscription that qualifies a person as a member
- **Service (Servizio)**: Physical resource offered by CNC (box, rack, boat slip, etc.)
- **Service Rental (Affitto Servizio)**: Annual contract between member and CNC for use of a service
- **Renewal Right (Diritto di Rinnovo)**: Member's right to keep the same service every year
- **Preemption Right (Diritto di Prelazione)**: Priority of existing member over new members for a service
- **Waiting List (Lista di Attesa)**: Queue of requests for unavailable services
- **Exclusion (Esclusione)**: Definitive removal of member decided by the board
- **Box Discount (Sconto Box)**: Price reduction on some services when also renting a Box

---

## Core Domain: Services Management

### Aggregates

#### 1. Member (Aggregate Root)
```
Member
├── MemberId (Identity)
├── PersonalInfo
│   ├── FirstName
│   ├── LastName
│   ├── Email
│   ├── Phone
│   └── Address
├── MembershipCard
│   ├── StartDate
│   ├── ExpirationDate
│   ├── PaymentStatus (Paid, Unpaid)
│   └── Amount
├── MemberStatus (Active, NotRenewed, Excluded)
├── ExclusionDate?
└── RegistrationDate

Behaviors:
+ RenewMembershipCard()
+ MarkCardPayment()
+ Exclude(deliberationDate)
+ HasRenewalRight() -> bool
+ HasPreemptionRight() -> bool
```

#### 2. Service (Aggregate Root)
```
Service
├── ServiceId (Identity)
├── ServiceType (Enum)
│   ├── OpenBoardRack
│   ├── DinghyYardBoatSlips
│   ├── Box
│   ├── ClosedBoardRack
│   ├── LargeWoodenLockers
│   ├── StandardLockers
│   ├── SurfStorageWorkshopSide
│   ├── ClosedSUPStorage
│   ├── ExternalCanoeRack
│   ├── VentenaRiverLandSlips
│   └── TavolloRiverBoatSlips
├── Number? (for numbered services)
├── ServiceCategory
│   ├── NumberedService (fixed quantity)
│   └── VariableService (river slips, depends on boats)
├── Availability
│   ├── TotalCapacity? (for numbered services)
│   ├── OccupiedSlots
│   └── IsAvailable -> bool
├── BasePrice
├── HasBoxDiscount -> bool
└── SizeRequirements? (for river boat slips)

Behaviors:
+ CalculateAvailability() -> int
+ Occupy()
+ Release()
+ CalculatePriceForMember(member) -> decimal
```

#### 3. ServiceRental (Aggregate Root)
```
ServiceRental
├── RentalId (Identity)
├── MemberId (Foreign Reference)
├── ServiceId (Foreign Reference)
├── ReferenceYear
├── StartDate
├── ExpirationDate
├── AnnualAmount
├── PaymentStatus (Paid, Unpaid, Late)
├── PaymentDate?
├── AssignedNumber? (e.g. Box #15)
├── IsContinuationOfPrevious (for renewal right)
└── RequestDate

Behaviors:
+ MarkAsPaid(paymentDate)
+ Renew(newYear) -> ServiceRental
+ IsLate(referenceDate) -> bool
+ CalculateAmount(service, hasBox) -> decimal
+ Terminate()
```

#### 4. WaitingList (Aggregate Root)
```
WaitingList
├── WaitingListId (Identity)
├── ServiceType
├── PendingRequests: List<WaitingRequest>
│   ├── RequestId
│   ├── MemberId
│   ├── RequestDate
│   ├── Priority (calculated based on member seniority)
│   └── Notes?
└── Ordered queue by priority/date

Behaviors:
+ AddRequest(memberId, requestDate)
+ RemoveRequest(requestId)
+ GetNextInQueue() -> MemberId?
+ NotifyAvailability(serviceId)
```

---

## Domain Services

### ServiceAvailabilityService
```
+ CheckAvailability(serviceType) -> bool
+ CalculateAvailableSlots(serviceType) -> int
+ AssignService(memberId, serviceId) -> Result
+ ManagePreemption(serviceId, requestingMember, currentMembers) -> Result
```

### RenewalService
```
+ RenewAllServices(memberId, newYear) -> List<ServiceRental>
+ VerifyRenewalRight(memberId, serviceId) -> bool
+ ApplyRenewalRight(serviceRental) -> ServiceRental
```

### PaymentService
```
+ CalculateTotalAnnualAmount(memberId, year) -> decimal
+ RecordMembershipPayment(memberId, date)
+ RecordServicePayment(rentalId, date)
+ GetLateMembers(referenceDate) -> List<Member>
```

---

## Value Objects

### PersonalInfo
```
{
  FirstName: string
  LastName: string
  Email: email
  Phone: string
  Address: string
}
```

### ValidityPeriod
```
{
  StartDate: date
  ExpirationDate: date

  IsValid(referenceDate) -> bool
  IsExpired(referenceDate) -> bool
}
```

### PaymentAmount
```
{
  Amount: decimal
  Currency: string (default EUR)

  ApplyDiscount(percentage) -> PaymentAmount
}
```

### BoatDimensions
```
{
  Length: decimal (meters)
  Width: decimal (meters)
  Category: enum (Small, Medium, Large)
}
```

---

## Enums

### MemberStatus
- Active
- NotRenewed (expired card, awaiting board decision)
- Excluded (board deliberation)

### PaymentStatus
- Paid
- Unpaid
- Late (past due date)

### ServiceType
- OpenBoardRack
- DinghyYardBoatSlips
- Box
- ClosedBoardRack
- LargeWoodenLockers
- StandardLockers
- SurfStorageWorkshopSide
- ClosedSUPStorage
- ExternalCanoeRack
- VentenaRiverLandSlips
- TavolloRiverBoatSlips

---

## Domain Events

### MemberEvents
- `MemberRegistered`
- `MembershipCardRenewed`
- `MembershipCardPaid`
- `MemberExcluded`

### ServiceEvents
- `ServiceRented`
- `ServiceReleased`
- `ServiceRenewed`
- `ServicePaymentRecorded`
- `ServiceAvailable` (trigger for waiting list)

### WaitingListEvents
- `RequestAddedToWaitingList`
- `RequestRemovedFromWaitingList`
- `AvailabilityNotificationSent`

---

## Repositories

### IMemberRepository
```
+ FindById(memberId) -> Member
+ FindAll() -> List<Member>
+ FindByStatus(memberStatus) -> List<Member>
+ FindMembersWithUnpaidCard() -> List<Member>
+ Save(member)
+ Delete(memberId)
```

### IServiceRepository
```
+ FindById(serviceId) -> Service
+ FindByType(serviceType) -> List<Service>
+ FindAvailable() -> List<Service>
+ FindByNumber(type, number) -> Service?
+ Save(service)
```

### IServiceRentalRepository
```
+ FindById(rentalId) -> ServiceRental
+ FindByMember(memberId, year?) -> List<ServiceRental>
+ FindByService(serviceId, year?) -> List<ServiceRental>
+ FindUnpaid(year?) -> List<ServiceRental>
+ FindLate(referenceDate) -> List<ServiceRental>
+ FindByYear(year) -> List<ServiceRental>
+ Save(rental)
+ Delete(rentalId)
```

### IWaitingListRepository
```
+ FindByServiceType(type) -> WaitingList
+ Save(waitingList)
```

---

## Application Services (Use Cases)

### MemberApplicationService
```
+ RegisterNewMember(personalInfo) -> MemberId
+ RenewMembershipCard(memberId) -> Result
+ RecordCardPayment(memberId, date) -> Result
+ ExcludeMember(memberId, deliberationDate) -> Result
+ GetMembersList() -> List<MemberDTO>
+ GetMembersWithUnpaidCard() -> List<MemberDTO>
+ GetMemberDetails(memberId) -> MemberDetailsDTO
```

### ServiceApplicationService
```
+ GetAvailableServices() -> List<ServiceDTO>
+ GetMemberRentedServices(memberId) -> List<ServiceRentalDTO>
+ RequestServiceRental(memberId, serviceType, number?) -> Result
+ ReleaseService(rentalId) -> Result
+ RecordServicePayment(rentalId, date) -> Result
+ RenewServices(memberId, newYear) -> Result
+ GetUnpaidServices() -> List<ServiceRentalDTO>
+ GetLateServices() -> List<ServiceRentalDTO>
```

### WaitingListApplicationService
```
+ AddRequestToWaitingList(memberId, serviceType) -> Result
+ RemoveFromWaitingList(requestId) -> Result
+ GetWaitingList(serviceType) -> List<WaitingRequestDTO>
+ NotifyNextInQueue(serviceId) -> Result
```

---

## Business Rules

### R1: Renewal Right
A member who rented a service in year N has the automatic right to renew it for year N+1, keeping the same number (if applicable).

**Implementation**:
- `ServiceRental.Renew()` verifies `IsContinuationOfPrevious`
- `RenewalService.VerifyRenewalRight()` checks rental history

### R2: Preemption Right
When a service becomes available, existing members who previously occupied it have priority over new members.

**Implementation**:
- `ServiceAvailabilityService.ManagePreemption()` orders requesters by seniority
- `WaitingList` calculates priority based on member history

### R3: Box Discount
The services ClosedBoardRack, SurfStorageWorkshopSide, and ExternalCanoeRack get a discount if the member also rents a Box.

**Implementation**:
- `ServiceRental.CalculateAmount()` verifies presence of active Box
- `Service.CalculatePriceForMember()` applies discount logic

### R4: Non-Paying Member
A member who doesn't pay the membership card remains a member until exclusion by the board.

**Implementation**:
- `Member.MemberStatus` distinguishes between Active/NotRenewed/Excluded
- Only `Member.Exclude()` with board deliberation changes status definitively

### R5: Variable Availability
River slips (Ventena, Tavollo) have variable availability based on number and size of boats present.

**Implementation**:
- `Service.CalculateAvailability()` for variable services uses empirical logic
- `ServiceAvailabilityService` allows manual capacity updates

### R6: FIFO Waiting List with Priority
When a service becomes available, the first member in the waiting list is notified (considering preemption right priority).

**Implementation**:
- `WaitingList` maintains ordered queue
- `ServiceAvailable` event triggers automatic notification

---

## Extension: Courses Context (Future)

### Aggregates

#### Course
```
Course
├── CourseId
├── CourseType (Sailing, Windsurf, SUP, Diving, WingFoil, KiteSurf)
├── ClientType (Private, MunicipalSummerCamp)
├── Period (StartDate, EndDate)
├── Participants: List<Participant>
├── TotalPrice
├── PaymentStatus
└── Notes

Behaviors:
+ AddParticipant()
+ RecordPayment()
+ GenerateReport()
```

#### Participant
```
Participant
├── FirstName
├── LastName
├── Age
├── EmergencyContact
└── MedicalCertificate?
```

### Application Service
```
CourseApplicationService
+ CreateCourse(courseType, clientType, period) -> CourseId
+ AddParticipant(courseId, participant) -> Result
+ RecordCoursePayment(courseId, amount, date) -> Result
+ GeneratePaymentReport(period) -> ReportDTO
+ GetActiveCourses() -> List<CourseDTO>
```

---

## Database Schema Hints

### Tables
- `Members` (aggregate root)
- `MembershipCards` (owned by Member)
- `Services` (aggregate root)
- `ServiceRentals` (aggregate root, references Member + Service)
- `WaitingLists` (aggregate root)
- `WaitingRequests` (owned by WaitingList)
- `Events` (event sourcing optional)

### Key Relationships
- `ServiceRentals.MemberId` → `Members.Id`
- `ServiceRentals.ServiceId` → `Services.Id`
- `WaitingRequests.MemberId` → `Members.Id`
- `WaitingRequests.WaitingListId` → `WaitingLists.Id`

---

## Architecture Recommendations

### Layering
```
┌─────────────────────────────────┐
│   Presentation Layer            │
│   (API, Web UI, Reports)        │
└─────────────────────────────────┘
           ↓
┌─────────────────────────────────┐
│   Application Layer             │
│   (Use Cases, DTOs)             │
└─────────────────────────────────┘
           ↓
┌─────────────────────────────────┐
│   Domain Layer                  │
│   (Aggregates, Services, Events)│
└─────────────────────────────────┘
           ↓
┌─────────────────────────────────┐
│   Infrastructure Layer          │
│   (Repositories, DB, Email)     │
└─────────────────────────────────┘
```

### Technology Suggestions
- **Backend**: .NET 8+ / Java Spring / Node.js TypeScript
- **Database**: PostgreSQL / SQL Server
- **ORM**: Entity Framework Core / Hibernate / TypeORM
- **API**: REST + OpenAPI / GraphQL
- **Authentication**: JWT / OAuth2
- **Reports**: Crystal Reports / SSRS / Custom PDF generation

### Key Patterns
- **Repository Pattern**: Data access abstraction
- **Unit of Work**: Transaction management
- **Domain Events**: Cross-aggregate communication
- **Specification Pattern**: Complex queries (waiting list priority)
- **Strategy Pattern**: Service price calculation with discount
