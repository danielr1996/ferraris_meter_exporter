# ferraris_meter_exporter

prometheus exporter for [ferraris meter / eletricity counter](https://en.wikipedia.org/wiki/Electricity_meter#Electromechanical)

> This software requires additional hardware. See the Hardware section for additional info.
> Additionally it is designed to run on a Raspberry Pi or similiar devices that supports GPIO.
> Specifically it uses the [gpiod librariy](https://github.com/warthog618/gpiod) to read the GPIO Pins

# What is measured?
The `ferraris_meter_rpms` metric indicates the total number of rotations.

Due to floating point arithemic inaccuracies the total power usage in kWh is not reported,
rather divide the rpms by 75 (depending on your ferraris meter, see the Hardware section for details) in prometheus:
```
ferraris_meter_rpms/75
```
This reduces the inaccuracy to 1 floating point division overall instead of 1 floating point division each rotation
which would get more inaccurate over time the more rotations are counted.

# Run
> ferraris_meter_exporter assumes that the IR Sensor is wired to GPIO2 on the Raspberry Pi and that
> the GPIO can be accessed via `/dev/gpiochip0`. See the Build section to configure the values yourself.

To run simply download the binary from the releases page and use

```
./ferraris_meter_exporter
```

Your metrics are now served at `http://raspberrypi:2112/metrics`

# Build
## Configuration
If you want to wire the IR Sensor to a pin different to GPIO2, change the gpio device or change the port change
the following values in `src/main.go`

Refer to the [gpiod documentation](https://github.com/warthog618/gpiod#chip-initialization) for available
`CHIP` values.

```
const PIN = rpi.GPIO2
const PORT = "2112"
const CHIP = "gpiochip0"
```

To build directly on the Raspberry Pi (might be slow) run
```
go build -o ferraris_meter_exporter .
```

To cross compile for 32 Bit Raspberry Pi (all Raspbian/RaspiOS images) on your PC run
```
GOOS=linux GOARCH=arm go build -o ferraris_meter_exporter .
```

To cross compile for 64 Bit Raspberry Pi (Ubuntu) on your PC (not tested as I don't use Ubuntu) run
```
GOOS=linux GOARCH=arm64 go build -o ferraris_meter_exporter .
```

Now run on the Raspberry Pi with

```
chmod a+x ferraris_meter_exporter
./ferraris_meter_exporter
```

# Hardware
This project is inspired by [SchimmerMediaHD](https://www.youtube.com/watch?v=ZZkQDy53GcM)(german)
and [mwinkler](https://mwinkler.jimdo.com/smarthome/aktoren-sensoren/stromz%C3%A4hler/)(german).

For a full guide on the hardware setup refer to the links. I will just summarize them and outline the differences.

## How does the measuring work?
A ferraris meter has a rotating disc that indicates the energy usage. The ferraris meter should have a label that states something like 
`X U/kWh` (U=Umdrehung=rpms) that means each rotation 1/X kWh is consumed. The disc should also have a red or black marker that we can use to
detect when a rotation has completed.

## TCRT5000
To detect a rotation I used the [TCRT500 IR sensor](https://www.amazon.de/AZDelivery-TRCT5000-Infrarot-Hindernis-Vermeidung/dp/B07DRCKV3X/ref=sr_1_1_sspa?__mk_de_DE=%C3%85M%C3%85%C5%BD%C3%95%C3%91&crid=1R423XVRM8SQY&dchild=1&keywords=tcrt5000&qid=1600905822&sprefix=tcrt500%2Caps%2C157&sr=8-1-spons&psc=1&spLa=ZW5jcnlwdGVkUXVhbGlmaWVyPUFHTlRCMjZJTEhHMFQmZW5jcnlwdGVkSWQ9QTA1Njg1NzAyU0tNRlFLSExGQjM2JmVuY3J5cHRlZEFkSWQ9QTA3MTcyNDkxUlVGQ1JDMUE2NEowJndpZGdldE5hbWU9c3BfYXRmJmFjdGlvbj1jbGlja1JlZGlyZWN0JmRvTm90TG9nQ2xpY2s9dHJ1ZQ==),
but every sensor that has a digital output will work.

The TCRT5000 has 4 outputs: A0, D0, VCC, GND. 
1. Wire VCC to 3.3V on the Raspberry Pi (don't use 5V, it might break your TCRT5000)
2. Wire GND to GND on the Raspberry Pi
3. Wire D0 to a GPIO on the Raspberry Pi (use GPIO2 if you want to use the prebuilt binary)
4. A0 will be unused as we don't use the analog signal