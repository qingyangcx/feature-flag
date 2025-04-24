#      Feature flag

## About

Feature flags (also known as feature toggles or feature switches) are a software development and product experimentation technique that turns certain functionality on and off during runtime, without deploying new code.

This project has stood the test of time in a demanding production environment, handling peak loads exceeding 3000 queries per second. While cloud-based feature flag services like [Optimizely](https://www.optimizely.com/optimization-glossary/feature-flags/) offer robust solutions, building your own can provide tailored functionality to meet specific company needs. This project offers a hands-on opportunity to explore the core concepts of feature flags and potentially develop a customized solution.



## Feature-flag use cases and benefits

1. Testing in production
2. Canary releases
3. Quicker release cycles
4. Rollback / kill switch
5. Server-side A/B Testing

## Fast start

1. Install a mysql and launch
2. Init database by execute sql/init.sql. Fill the database config to conf/test.json
3. Run feature-flag. 
`go run cmd/main.go --work_dir="$(pwd)" --config_file="conf/test.json" `
4. Create first feature-flag

```
curl --location '127.0.0.1:8000/feature' \
--header 'Content-Type: application/json' \
--data '{
    "name":"feature_demo",
    "key": "feature_demo_key",
    "blacklist": [],
    "namespace": "demo",
    "valid":true,
    "feature_values": [
        {
            "value": "value1",
            "whitelist":[],
            "traffic": 600,
            "default":true
        },
        {
            "value": "value2",
            "whitelist":[],
            "traffic": 400
        }
    ]
}'
```

5. Try to split, the feature_id should be the id of feature-flag that you just created
```
curl --location --request GET '127.0.0.1:8000/split' \
--header 'Content-Type: application/json' \
--data '{
    "namespace":"demo",
    "identity":"test",
    "feature_id":21
}'
```

You can get to know all interfaces by looking at the file main.go 