# Reproducing my cloud-run issue

https://stackoverflow.com/questions/57506187/cloud-run-process-fails-with-500-status-code-and-a-membarrier-gvisor-error/57518449#57518449

## Introduction

My service hangs on large transfers and I don't understand why. It works perfectly fine on my Macbook or in a local container.

This project aims to reproduce the problem with as few lines as possible.

# About the code

No rocket science. We set the headers and then manually write "A" to the `http.ResponseWriter`. By default it fetches 200M. Size can be adjusted with a query-parameter.

`https://servicelink/?size=100` will download a 100M file.

I used the server and tracing packages from gocloud.dev as I did with my other service.

# If you want to build and run

You can skip this and just follow the next step from our service. This roughly follows the cloud-run quickstart. Note: I will remove the service in a few days.

```
gcloud builds submit --tag eu.gcr.io/redia-internal/reproduce-cloud-run-bug
gcloud beta run deploy --image eu.gcr.io/redia-internal/reproduce-cloud-run-bug --platform managed
```

# Reproducing

Then from your terminal enter the following

```sh
# 1 MB data works fine
curl https://reproduce-cloud-run-bug-acyeuzbsfq-ew.a.run.app?size=1 > wtf.dat && ls -alh wtf.dat
-rw-r--r--  1 andri  staff   1.0M Aug 16 11:13 wtf.dat
# 100 MB data fails
curl https://reproduce-cloud-run-bug-acyeuzbsfq-ew.a.run.app?size=100 > wtf.dat && ls -alh wtf.dat
cat wtf.dat

<html><head>
<meta http-equiv="content-type" content="text/html;charset=utf-8">
<title>500 Server Error</title>
</head>
<body text=#000000 bgcolor=#ffffff>
<h1>Error: Server Error</h1>
<h2>The server encountered an error and could not complete your request.<p>Please try again in 30 seconds.</h2>
<h2></h2>
</body></html>
```

# Things to rule out

- I'm using the gocloud.dev server package, that might be causing this.
