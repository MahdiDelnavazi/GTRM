-- example HTTP POST script which demonstrates setting the
-- HTTP method, body, and adding a header

wrk.method = "PUT"
wrk.body = '{}'
wrk.headers["Content-Type"] = "application/json"
