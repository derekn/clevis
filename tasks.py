import json
from os import mkdir
from pathlib import Path
from shutil import rmtree

from invoke import task


def get_version(c):
	return json.loads(c.run("gitversion -output json", hide=True).stdout)['SemVer']

@task
def clean(c):
	"""clean dist dir"""

	rmtree('dist', ignore_errors=True)
	mkdir('dist')

@task
def build(c, version=None, target=None):
	"""build release version"""

	version = version if version else get_version(c)
	output = f"dist/clevis-{version}"
	env = {}

	if target:
		env['GOOS'], env['GOARCH'] = target.split('/', 1)
		output += f"-{env['GOOS']}-{env['GOARCH']}"

	c.run(f"go build -ldflags '-s -w -X main.version={version}' -o '{output}'", echo=True, env=env)

@task(clean)
def release(c, targets=None, compress=True):
	"""do release builds and compress"""

	version = get_version(c)
	targets = targets if targets else ['darwin/amd64', 'darwin/arm64', 'linux/amd64', 'linux/arm', 'linux/arm64']

	for target in targets:
		build(c, version, target)

	if compress:
		files = Path('dist').glob('*')
		for file in files:
			c.run(f"xz --compress --threads 0 {file}", echo=True)
