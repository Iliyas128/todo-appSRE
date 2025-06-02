# Postmortem: Database Downtime Simulation

## Summary
On June 2, 2025, the ToDo application experienced an artificial database failure that caused the application to return 500 errors for all POST requests.

## Impact
All users attempting to create tasks received a 500 error.

## Root Cause
The database service was stopped manually to simulate unavailability.

## Timeline
- 12:00 - Application running normally
- 12:05 - Database container stopped
- 12:06 - Alerts triggered in Prometheus
- 12:10 - Database restarted
- 12:11 - Service recovered

## Resolution
Restarted the database and validated app health.

## Action Items
- Implement retry logic in DB layer
- Add health check for DB connectivity
- Setup alert escalation

## Lessons Learned
Simulations are crucial for training incident response and validating observability.
