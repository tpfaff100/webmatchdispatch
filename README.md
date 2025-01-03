# webmatchdispatch
Very simple GOLANG demonstration of dynamic https string->function dispatch table set up to run methods by way of a Dictionary/Hashtable list
<pre>
Usage:

------------------
Using OPENSSL, generate unsafe https key so we can run an https server easily:
openssl genrsa 2048 > server.key
chmod 400 server.key
openssl req -new -x509 -nodes -sha256 -days 365 -key server.key -out server.crt
-------------------
go build *.go
./dictdispatchs  

Then load the web browser up and try two links:

https://localhost:9876/login
https://localhost:9876

This example shows how to setup function dispatch tables in golang
Easy peasy.  I just haven't seen it done like this anywhere yet.

I have used very complex tables for other applications in the past I like how the function names all go at the top of the main() file and are easy to find, rather than digging through the main() procedure, looking for each tag->function match.

Almost forgot- the .conf file determines what PORT is used for the https connection.
</pre>
