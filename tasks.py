from os import mkdir
from pathlib import Path
from shutil import rmtree

from invoke import task


@task
def build(c):
	targets = ['darwin/amd64', 'darwin/arm64', 'linux/amd64', 'linux/arm', 'linux/arm64']

	rmtree('dist', ignore_errors=True)
	mkdir('dist')

	for target in targets:
		os, arch = target.split('/')
		c.run(f"go build -ldflags '-s -w' -o 'dist/clevis-{os}-{arch}'", echo=True, env={'GOOS': os, 'GOARCH': arch})

@task(build)
def release(c):
	files = Path('dist').glob('*')

	for file in files:
		c.run(f"xz --compress --threads 0 {file}", echo=True)
