import http from 'k6/http';
import { check, fail } from 'k6';

export  function callGraphQL(url, payload, token, variables) {
    const headers = {
        'Content-Type': 'application/json',
    };

    if (token) {
        headers['Authorization'] = 'Bearer ' + token;
    }

    const res = http.post(url, JSON.stringify({ query: payload, variables }), { headers });

    if (res.status === 200) {
        let body = JSON.parse(res.body);

        const checkOutput = check(
            body,
            {
                'no graphql errors': (body) => body === null || body.errors === null || body.errors === undefined,
            }
        );
        if (!checkOutput) {
            console.error("Response: " + JSON.stringify(res.body))

            let error = body.errors[0]
            const msg = JSON.stringify(error.message)
            console.error("Error: " + JSON.stringify(error))
            fail(msg)
            return null
        }
        return body.data
    }

    console.log(res.body)
    fail(`Status of response is wrong`)
}


export function uploadFile(url, payload, token, variables, fileBinary, variablePath) {
    let headers = {
    };

    if (token !== "") {
        headers['Authorization'] = 'Bearer ' + token;
    }

    const data = {
        operations: JSON.stringify({ query: payload, variables }),
        map: JSON.stringify({ "0": [variablePath] }),
        "0": http.file(fileBinary, 'test.jpg', 'image/jpeg'),
    };

    let res = http.post(url, data, { headers: headers });

    if (res.status === 200) {
        let body = JSON.parse(res.body);

        const checkOutput = check(
            body,
            {
                'no graphql errors': (body) => body === null || body.errors === null || body.errors === undefined,
            }
        );
        if (!checkOutput) {
            console.error("Response: " + JSON.stringify(res.body))

            let error = body.errors[0]
            const msg = JSON.stringify(error.message)
            console.error("Error: " + JSON.stringify(error))
            fail(msg)
            return null
        }
        return body.data
    }

    fail(`Status of response is wrong`)
}
