[![Actions Status](https://github.com/jfbramlett/go-template/workflows/Go/badge.svg)](https://github.com/jfbramlett/go-template/actions)

# Start DB container
start docker container d8b924fe91fe

# Sample Table Test Results
- 1,000,000 rows in sample table
- sample query with in clause containing 100 random UUID's
- 10 test runs
- avg query response: 2ms


# Metadata Table Test Results
- 4,000,000 rows in metadata table (1,000,000 assets * 4 attributes per)
- sample query with in clause containing 100 random UUID's
- 10 test runs
- avg query response: 6ms (max was 9ms, min was 5ms)



