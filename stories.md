# TIP JAR USER STORIES

## Iteration 1: Core MVP - "Get It Working"

### Goal: Basic functional tip jar with simple offense tracking

- [x] **As a** new user, **I want to** sign up using OAuth (Google/GitHub/Discord), **so that** I can quickly create an account without another password
- [x] **As a** returning user, **I want to** log in with my OAuth provider, **so that** I can access my tip jars
- [x] **As a** user, **I want to** create a new tip jar with a name and description, **so that** I can start tracking offenses with my group
- [x] **As a** user, **I want to** join an existing tip jar using an invite code, **so that** I can participate with my friends
- [x] **As a** jar admin, **I want to** generate and share an invite code, **so that** others can join my tip jar
- [x] **As a** user, **I want to** see a list of all tip jars I'm a member of, **so that** I can switch between different groups
- [ ] **As a** user, **I want to** report an offense committed by another member, **so that** it's tracked in the ledger
- [ ] **As an** offender, **I want to** mark a monetary offense as paid, **so that** my debt is cleared
- [ ] **As a** user, **I want to** see all pending offenses in a tip jar, **so that** I know what's currently owed
- [ ] **As a** user, **I want to** see each member's current balance in the jar, **so that** I know who owes what

## Iteration 2: Custom Offenses - "Make It Personalized"

### Goal: Let jars define their own offense types with different cost types

- [ ] **As a** jar admin, **I want to** create custom offense types with specific costs, **so that** our jar has relevant penalties
- [ ] **As a** jar admin, **I want to** set non-monetary costs (actions/items/services), **so that** penalties can be creative
- [ ] **As a** user, **I want to** see all available offense types in my jar, **so that** I know what can be reported
- [ ] **As a** user, **I want to** select from predefined offense types when reporting, **so that** reporting is consistent
- [ ] **As a** user, **I want to** add notes when reporting an offense, **so that** I can provide context
- [ ] **As a** user, **I want to** override the default cost when reporting, **so that** I can adjust for specific situations
- [ ] **As a** jar admin, **I want to** edit existing offense types, **so that** I can adjust penalties as needed
- [ ] **As a** jar admin, **I want to** deactivate offense types without deleting them, **so that** I can preserve history

## Iteration 3: Payment Verification - "Make It Trustworthy"

### Goal: Add proof and verification for payments, especially non-monetary ones

- [ ] **As an** offender, **I want to** upload proof of completing an action-based penalty, **so that** others can verify
- [ ] **As a** user, **I want to** attach receipts/photos to payments, **so that** there's proof of settlement
- [ ] **As a** witness, **I want to** verify someone's action-based payment, **so that** settlements are legitimate
- [ ] **As a** user, **I want to** see which specific offenses a payment covered, **so that** the ledger is transparent
- [ ] **As an** offender, **I want to** batch pay multiple offenses at once, **so that** I can clear my debts efficiently
- [ ] **As a** user, **I want to** see all settled offenses in a tip jar, **so that** I can review payment history

## Iteration 4: Personal Dashboard - "Make It Personal"

### Goal: Give users a cross-jar view of their obligations and history

- [ ] **As a** user, **I want to** see and edit my profile information, **so that** others can identify me correctly
- [ ] **As a** user, **I want to** see all my pending offenses across all jars, **so that** I know my total obligations
- [ ] **As a** user, **I want to** see my payment history across all jars, **so that** I can track what I've settled
- [ ] **As a** user, **I want to** see my offenses grouped by jar, **so that** I can prioritize which debts to pay
- [ ] **As a** user, **I want to** see my total balance across all jars, **so that** I know my overall standing
- [ ] **As a** user, **I want to** see offenses I've reported across all jars, **so that** I can track my reporting activity

## Iteration 5: Timeline & Filtering - "Make It Browsable"

### Goal: Add chronological views and filtering for better ledger navigation

- [ ] **As a** user, **I want to** see a chronological timeline of all jar activity, **so that** I can track what happened when
- [ ] **As a** user, **I want to** filter the jar timeline by date range, **so that** I can focus on specific periods
- [ ] **As a** user, **I want to** filter jar activity by member, **so that** I can see individual participation

## Iteration 6: Disputes & Moderation - "Make It Fair"

### Goal: Add dispute resolution and moderation tools

- [ ] **As an** offender, **I want to** acknowledge an offense reported against me, **so that** I accept responsibility
- [ ] **As an** offender, **I want to** dispute an offense reported against me, **so that** unfair penalties can be reviewed
- [ ] **As** jar members, **I want to** vote on disputed offenses, **so that** the group can decide fairly
- [ ] **As a** jar admin, **I want to** reverse/forgive offenses, **so that** mistakes can be corrected
- [ ] **As a** jar member, **I want to** suggest new offense types, **so that** the jar stays relevant

## Iteration 7: Notifications - "Make It Engaging"

### Goal: Keep users informed about jar activity

- [ ] **As an** offender, **I want to** be notified when someone reports an offense against me, **so that** I'm aware immediately
- [ ] **As a** reporter, **I want to** be notified when an offense I reported is paid, **so that** I know it's settled
- [ ] **As a** jar member, **I want to** see notifications for all jar activity, **so that** I stay informed
- [ ] **As a** user, **I want to** configure which notifications I receive, **so that** I'm not overwhelmed

## Iteration 8: Admin & Analytics - "Make It Powerful"

### Goal: Add administrative tools and data insights

- [ ] **As a** user, **I want to** leave a tip jar, **so that** I'm no longer part of that group
- [ ] **As a** jar admin, **I want to** remove members from my tip jar, **so that** I can manage who participates
- [ ] **As a** jar admin, **I want to** set jar-wide rules about offense limits, **so that** things don't get out of hand
- [ ] **As a** jar admin, **I want to** export the jar's ledger as CSV, **so that** I can analyze it externally
- [ ] **As a** user, **I want to** export my personal offense history, **so that** I have records for myself
- [ ] **As a** jar admin, **I want to** see analytics about offense patterns, **so that** I can understand group dynamics
