# MBX-IOT

This is an IoT project to monitor for mail delivery and watch for cars passing by on my street.

## Doc

This repo uses [mkdocs](https://www.mkdocs.org/) ([help](https://mkdocs.readthedocs.io/en/0.10/)) and [github pages](https://help.github.com/articles/configuring-a-publishing-source-for-github-pages/) to host content at:

[https://tonygilkerson.github.io/mbx-iot/](https://tonygilkerson.github.io/mbx-iot/)

**Develop:**

```sh
pip3 install mkdocs-same-dir

mkdocs serve
# Edit content and review changes here:
open http://127.0.0.1:8000/
```

## LORA bug workaround

See this [issue](https://github.com/tonygilkerson/mbx-iot/issues/5) for more detail

* Install TinyGo the official way with brew
* Get everything working in vscode
  * Install the TinyGo plugin
  * Make sure the TinyGo target is set and `machine` is recognized in the IDE
* Modify the file as described below

```sh
# Make sure you have `pico` selected as the TinyGo target
# and the vscode is working then run the following
code  $(jq -r  '.["go.toolsEnvVars"].GOROOT' .vscode/settings.json)/src/machine/machine_rp2040_spi.go 

# Comment out the two occurrences of the following and save
	for spi.isBusy() {
		gosched()
	}
```
