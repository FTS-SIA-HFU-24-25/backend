# Backend

This is a UDP & TCP Server that receives messages as in Byte Array and sends back with the same data format.

## Receiving Data

For receiving data we specify them in multiple cases

### Establishing Connection

For a hardware to establish a connection with the server, first you need send a TCP request to the server and the server will response you with a UUID that you can use for your UDP data streams

It will to this following endpoint:

> http://<YOUR_DOMAIN>/api/connection


### Streaming UDP data from Serial

To stream your serial data that you received from your sensors, you need to follow this format!

#### Connection with specified length

**Connection beginn**

This is when you say to the server: `Hey, I'm this, I want to send this length of data`

> [UUID, VERSION, LENGTH, DATA_TYPE]

After that you send the data in chunks that you like, until you reached the length, that you have provided

> [UUID, LENGTH, DATA]

#### Streaming without knowing the length

**Connection beginn**

This is when you say to the server: `Hey, I'm this, I want to send data`

> [UUID, VERSION, 0, DATA_TYPE]

To send data

> [UUID, LENGTH, DATA]

After that you just stream data, until your message look like this. This will mean, that you want to stop connection!

> [UUID, VERSION, 0]

#### Specifying the data type

In the previous section, we know how to stream data to the server from your serial. Now you have specify what kind of data you have sent to the server.

These are the supported sensors and hardware:

  - GPS Sensor on Arduino mkr gsm 1400
  - AD8232 ECG Sensor
  - LM-35 Temperature Sensor


**Specifying the data**

This is the normal data format that you will send to the server

> [UUID, VERSION, DATA]

But you can also specify the data type that you have sent to the server. By specifying the first byte of your data by replacing the `DATA_TYPE` for these avaiable bytes.

  > This is in binary

  - 0001: GPS Sensor on Arduino mkr gsm 1400
  - 0002: AD8232 ECG Sensor
  - 0003: LM-35 Temperature Sensor

  > In bytes

  - 1: GPS Sensor on Arduino mkr gsm 1400
  - 2: AD8232 ECG Sensor
  - 3: LM-35 Temperature Sensor
