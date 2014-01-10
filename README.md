# go-modbus

## About

This [Go](http://golang.org/) package provides [Modbus](http://en.wikipedia.org/wiki/Modbus) access for client (master) applications to communicate with server (slave) devices, over both [TCP/IP](http://www.modbus.org/docs/Modbus_Messaging_Implementation_Guide_V1_0b.pdf) and [Serial Line/RTU](http://www.modbus.org/docs/Modbus_over_serial_line_V1_02.pdf) frame protocols.

Note that in modbus terminology, _client_ refers to the __master__ application or device, and the _server_ is the __slave__ waiting to respond to instructions, as shown in this transaction diagram:

![Modbus Transaction](http://i.imgur.com/Vgsqrb2.png)

This code was originally forked from [lubia/modbus](https://github.com/lubia/modbus) and repositioned as a pure client (master) library for use by controller applications.

## References
- [Modbus Technical Specifications](http://www.modbus.org/specs.php)
- [Modbus Interface Tutorial](http://www.lammertbies.nl/comm/info/modbus.html)
- [Modbus TCP/IP Overview](http://www.rtaautomation.com/modbustcp/)
- [Modbus RTU Protocol Overview](http://www.rtaautomation.com/modbusrtu/)

## Acknowledgements
- [Lubia Yang](http://www.lubia.me) for the [original modbus code](https://github.com/lubia/modbus) in Go
- [l.lefebvre](http://source.perl.free.fr/) for his excellent [modbus client](https://github.com/sourceperl/MBclient) and [server (slave device simulator)](https://github.com/sourceperl/mbserverd) code repositories

