# object

Config contains the configuration for the datum server


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**refresh\_interval**|`integer`|RefreshInterval determines how often to reload the config<br/>||
|[**server**](#server)|`object`|Server settings for the echo server<br/>|yes|
|[**db**](#db)|`object`||yes|
|[**redis**](#redis)|`object`|||
|[**tracer**](#tracer)|`object`|||
|[**sessions**](#sessions)|`object`|||

**Additional Properties:** not allowed  
<a name="server"></a>
## server: object

Server settings for the echo server


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**debug**|`boolean`|Debug enables debug mode for the server<br/>|no|
|**dev**|`boolean`|Dev enables echo's dev mode options<br/>|no|
|**listen**|`string`|Listen sets the listen address to serve the echo server on<br/>|yes|
|**shutdown\_grace\_period**|`integer`|ShutdownGracePeriod sets the grace period for in flight requests before shutting down<br/>|no|
|**read\_timeout**|`integer`|ReadTimeout sets the maximum duration for reading the entire request including the body<br/>|no|
|**write\_timeout**|`integer`|WriteTimeout sets the maximum duration before timing out writes of the response<br/>|no|
|**idle\_timeout**|`integer`|IdleTimeout sets the maximum amount of time to wait for the next request when keep-alives are enabled<br/>|no|
|**read\_header\_timeout**|`integer`|ReadHeaderTimeout sets the amount of time allowed to read request headers<br/>|no|
|[**tls**](#servertls)|`object`|TLS settings for the server for secure connections<br/>|no|
|[**cors**](#servercors)|`object`|CORS settings for the server to allow cross origin requests<br/>|no|

**Additional Properties:** not allowed  
<a name="servertls"></a>
### server\.tls: object

TLS settings for the server for secure connections


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**enabled**|`boolean`|Enabled turns on TLS settings for the server<br/>||
|**cert\_file**|`string`|CertFile location for the TLS server<br/>||
|**cert\_key**|`string`|CertKey file location for the TLS server<br/>||
|**auto\_cert**|`boolean`|AutoCert generates the cert with letsencrypt, this does not work on localhost<br/>||

**Additional Properties:** not allowed  
<a name="servercors"></a>
### server\.cors: object

CORS settings for the server to allow cross origin requests


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**allow\_origins**](#servercorsallow_origins)|`string[]`|||
|**cookie\_insecure**|`boolean`|CookieInsecure allows CSRF cookie to be sent to servers that the browser considers<br/>unsecured. Useful for cases where the connection is secured via VPN rather than<br/>HTTPS directly.<br/>||

**Additional Properties:** not allowed  
<a name="servercorsallow_origins"></a>
#### server\.cors\.allow\_origins: array

**Items**

**Item Type:** `string`  
<a name="db"></a>
## db: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**debug**|`boolean`|debug enables printing the debug database logs<br/>|no|
|**databaseName**|`string`|the name of the database to use with otel tracing<br/>|no|
|**driverName**|`string`|sql driver name<br/>|no|
|**multiWrite**|`boolean`|enables writing to two databases simultaneously<br/>|no|
|**primaryDbSource**|`string`|dsn of the primary database<br/>|yes|
|**secondaryDbSource**|`string`|dsn of the secondary database if multi-write is enabled<br/>|no|
|**cacheTTL**|`integer`|cache results for subsequent requests<br/>|no|
|**runMigrations**|`boolean`|run migrations on startup<br/>|no|

**Additional Properties:** not allowed  
<a name="redis"></a>
## redis: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**enabled**|`boolean`|||
|**address**|`string`|||
|**name**|`string`|||
|**username**|`string`|||
|**password**|`string`|||
|**db**|`integer`|||
|**dialTimeout**|`integer`|||
|**readTimeout**|`integer`|||
|**writeTimeout**|`integer`|||
|**maxRetries**|`integer`|||
|**minIdleConns**|`integer`|||
|**maxIdleConns**|`integer`|||
|**maxActiveConns**|`integer`|||

**Additional Properties:** not allowed  
<a name="tracer"></a>
## tracer: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**enabled**|`boolean`|||
|**provider**|`string`|||
|**environment**|`string`|||
|[**stdout**](#tracerstdout)|`object`|||
|[**otlp**](#tracerotlp)|`object`|||

**Additional Properties:** not allowed  
<a name="tracerstdout"></a>
### tracer\.stdout: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**pretty**|`boolean`|||
|**disableTimestamp**|`boolean`|||

**Additional Properties:** not allowed  
<a name="tracerotlp"></a>
### tracer\.otlp: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**endpoint**|`string`|||
|**insecure**|`boolean`|||
|**certificate**|`string`|||
|[**headers**](#tracerotlpheaders)|`string[]`|||
|**compression**|`string`|||
|**timeout**|`integer`|||

**Additional Properties:** not allowed  
<a name="tracerotlpheaders"></a>
#### tracer\.otlp\.headers: array

**Items**

**Item Type:** `string`  
<a name="sessions"></a>
## sessions: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**signingKey**|`string`|||
|**encryptionKey**|`string`|||

**Additional Properties:** not allowed  

