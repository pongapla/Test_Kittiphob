# Queue Management System - Test Cases

Status: NOT RUN

| ID | Test Scenario | Expected Result |
|---|---|---|
| TC-01 | Request first ticket | A0 |
| TC-02 | Request next ticket after A0 | A1 |
| TC-03 | Request next ticket after A9 | B0 |
| TC-04 | Request next ticket after Z8 | Z9 |
| TC-05 | Request ticket after Z9 | Reject request |
| TC-06 | Reset queue | Display 00 |
| TC-07 | Request ticket after reset | A0 |
| TC-08 | 100 concurrent ticket requests | 100 unique tickets in the same cycle |
| TC-09 | Concurrent reset and issuance | Consistent serialized operations |
| TC-10 | Navigate back from ticket page | Return to reception page |
| TC-11 | Refresh ticket page | Preserve ticket number |
| TC-12 | Restart backend | Persist queue state |
| TC-13 | Request invalid ticket ID | Return 404 |
| TC-14 | Two consecutive resets | Next ticket is A0 |
| TC-15 | Request ticket after reset with old ticket ID | Old ticket remains identifiable by unique ID |

## Testing Tools
- Go Test: Unit and integration tests
- Go concurrency tests
- Playwright: Browser E2E tests

## Execution Results
Results will be recorded after implementation.

No test is considered PASS until it has been
executed and its actual result verified.
