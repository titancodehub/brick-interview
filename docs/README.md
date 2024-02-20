---
title: Open API Spec v0.1.0
language_tabs:
  - shell: curl
language_clients:
  - shell: shell
toc_footers: []
includes: []
search: true
highlight_theme: darkula
headingLevel: 2

---

<!-- Generator: Widdershins v4.0.1 -->

<h1 id="open-api-spec">Open API Spec v0.1.0</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

<h1 id="open-api-spec-default">Default</h1>

## get__validate-bank-accounts

> Code samples

```shell
# You can also use wget
curl -X GET /validate-bank-accounts?account_number=string&bank_code=BRI \
  -H 'Accept: application/json'

```

`GET /validate-bank-accounts`

Validate bank account

<h3 id="get__validate-bank-accounts-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|account_number|query|string|true|none|
|bank_code|query|string|true|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|bank_code|BRI|
|bank_code|BNI|
|bank_code|BCA|
|bank_code|MANDIRI|

> Example responses

> 200 Response

```json
{
  "account_number": "string",
  "account_name": "string",
  "bank_code": "string"
}
```

<h3 id="get__validate-bank-accounts-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Success validate bank account|[BankAccount](#schemabankaccount)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Bank account not found|[ErrorResponse](#schemaerrorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## post__disbursements

> Code samples

```shell
# You can also use wget
curl -X POST /disbursements \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /disbursements`

Disburse money to bank account

> Body parameter

```json
{
  "merchant_id": "string",
  "amount": 1,
  "account_number": "string",
  "bank_code": "BRI",
  "reference": "string"
}
```

<h3 id="post__disbursements-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|object|true|none|
|» merchant_id|body|string|true|none|
|» amount|body|integer|true|none|
|» account_number|body|string|true|none|
|» bank_code|body|string|true|none|
|» reference|body|string|true|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|» bank_code|BRI|
|» bank_code|BNI|
|» bank_code|BCA|
|» bank_code|MANDIRI|

> Example responses

> 200 Response

```json
{
  "id": "string",
  "merchant_id": "string",
  "reference": "string",
  "amount": 0,
  "status": "string",
  "type": "string",
  "metadata": {},
  "created": "2019-08-24T14:15:22Z",
  "updated": "2019-08-24T14:15:22Z"
}
```

<h3 id="post__disbursements-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Success validate bank account|[Transaction](#schematransaction)|

<aside class="success">
This operation does not require authentication
</aside>

## post__webhooks_disbursements

> Code samples

```shell
# You can also use wget
curl -X POST /webhooks/disbursements \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /webhooks/disbursements`

Receive notification from bank server

> Body parameter

```json
{
  "transaction_id": "string",
  "status": "SUCCESS"
}
```

<h3 id="post__webhooks_disbursements-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|object|true|none|
|» transaction_id|body|string|true|none|
|» status|body|string|true|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|» status|SUCCESS|
|» status|FAILED|

> Example responses

> 200 Response

```json
{
  "error": "string"
}
```

<h3 id="post__webhooks_disbursements-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Success validate bank account|[ErrorResponse](#schemaerrorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

# Schemas

<h2 id="tocS_BankAccount">BankAccount</h2>
<!-- backwards compatibility -->
<a id="schemabankaccount"></a>
<a id="schema_BankAccount"></a>
<a id="tocSbankaccount"></a>
<a id="tocsbankaccount"></a>

```json
{
  "account_number": "string",
  "account_name": "string",
  "bank_code": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|account_number|string|true|none|none|
|account_name|string|true|none|none|
|bank_code|string|false|none|none|

<h2 id="tocS_Transaction">Transaction</h2>
<!-- backwards compatibility -->
<a id="schematransaction"></a>
<a id="schema_Transaction"></a>
<a id="tocStransaction"></a>
<a id="tocstransaction"></a>

```json
{
  "id": "string",
  "merchant_id": "string",
  "reference": "string",
  "amount": 0,
  "status": "string",
  "type": "string",
  "metadata": {},
  "created": "2019-08-24T14:15:22Z",
  "updated": "2019-08-24T14:15:22Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|true|none|none|
|merchant_id|string|true|none|none|
|reference|string|true|none|none|
|amount|integer|true|none|none|
|status|string|true|none|none|
|type|string|true|none|none|
|metadata|object|false|none|none|
|created|string(date-time)|true|none|none|
|updated|string(date-time)|true|none|none|

<h2 id="tocS_WebhookResponse">WebhookResponse</h2>
<!-- backwards compatibility -->
<a id="schemawebhookresponse"></a>
<a id="schema_WebhookResponse"></a>
<a id="tocSwebhookresponse"></a>
<a id="tocswebhookresponse"></a>

```json
{
  "message": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|message|string|true|none|none|

<h2 id="tocS_ErrorResponse">ErrorResponse</h2>
<!-- backwards compatibility -->
<a id="schemaerrorresponse"></a>
<a id="schema_ErrorResponse"></a>
<a id="tocSerrorresponse"></a>
<a id="tocserrorresponse"></a>

```json
{
  "error": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|error|string|true|none|none|

