import {group, sleep} from 'k6';
import { Rate } from 'k6/metrics';
import {callGraphQL} from "./graphql.js";

const rate = new Rate('success_graphql_requests');

export const options = {
    "vus": 2000,
    "iterations": 10000,
    "tags": {
        "test_type": "load",
        "api_name": "GraphQL"
    },
    thresholds: {
        checks: ['rate==1'],
        http_req_failed: ['rate<0.02'],
        success_graphql_requests: ['rate>0.98'], // graphql success requests should be more than 98%
        group_duration: ['p(95)<3000']
    },
};


function main() {
    group('get-post', function () {
        try {
            const query = `{  
              posts(count:10){id,author{id, name}, title}
            }`
            callGraphQL('http://localhost:8181/graphql', query, '',{})
            rate.add(true)
        } catch (e) {
            rate.add(false)
            throw e
        }
    });
}

export function setup() {
}

export default function () {
    main()
}
export function teardown(data) {
    // teardown code
}
