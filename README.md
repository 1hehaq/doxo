<div align="center">
  <img src="https://raw.githubusercontent.com/1hehaq/doxo/refs/heads/main/avatar/doxo.jpeg" width="75" alt="logo"/>

<h5>
    
_A lightweight CLI tool that bridges your terminal with_ [Discord](discord.com) _through webhooks. Send your terminal output directly to Discord channels with simple flags like_ <kbd>-txt</kbd>, <kbd>-plain</kbd>, _and_ <kbd>-tts</kbd>!

</h5>
    
</div>

<br>
<br>
<br>

> [!Important]
> <kbd>doxo</kbd> _is still in improvement phase, so feel free to leave issues here_ [**`issues`**](https://github.com/1hehaq/doxo/issues)

<h1></h1>

<br>
<br>


<h3>
  
<kbd>↺</kbd> Install</h3>

```bash
go install github.com/1hehaq/doxo@latest
```

```yml
➜ ~ doxo -h

usage:
  doxo [flags] [message]

Flags:
   -config   : path to config file (default: ~/.doxo/doxo.json)
   -txt      : send output as text file
   -plain    : send as plain text message
   -tts      : send message with text-to-speech enabled
   -help     : show help message
```
<h1></h1>

<br>
<br>

<h3>
  
<kbd>↺</kbd> setup
  
</h3>

* _Create a Discord webhook in your server <br>_
* _Run_ <kbd>doxo</kbd> _to generate config file_ <br>
* _Add webhook URL to_ `~/.doxo/config.json` <br>

<h1></h1>

<br>
<br>

<h3>
  
<kbd>↺</kbd> using
  
</h3>

```bash
doxo -plain 'Target of month is: hackerone.com'
```

```bash
subfinder -d hackerone.com | doxo -txt && doxo -tts 'Subdomain Enumeration of hackerone.com completed!'
```

```bash
doxo -tts 'Hey hackers'
```

<h1></h1>

<br>
<br>

<kbd>↺</kbd> examples


|                                      |                                 |
| :----------------------------------: | :-----------------------------: |
|              <kbd>**Example**</kbd>         |  <kbd>**Usage**</kbd>|
|![webhook](https://github.com/user-attachments/assets/ef26a74b-7380-4d1d-bf5d-e2ff59222332)|![run](https://github.com/user-attachments/assets/2044a205-d360-4d96-84ec-69edddd97c88)|
|![recieve](https://github.com/user-attachments/assets/2d130750-f8b7-490e-bee6-9cedc2dd5c7f)|![send](https://github.com/user-attachments/assets/d8a7fbf0-81af-4447-a707-f9dc70331759)|
|![send](https://github.com/user-attachments/assets/b8bf7d34-c47e-4eb4-a72a-f70a37c5660f)|![recieve](https://github.com/user-attachments/assets/646c60e6-65cf-435f-9f96-1b5d189895be)


<br>
<br>
<br>

<h6 align="center">kindly for hackers</h6>


<div align="center">
  <a href="https://github.com/1hehaq"><img src="https://img.icons8.com/material-outlined/20/808080/github.png" alt="GitHub"></a>
  <a href="https://twitter.com/1hehaq"><img src="https://img.icons8.com/material-outlined/20/808080/twitter.png" alt="X"></a>
</div>
