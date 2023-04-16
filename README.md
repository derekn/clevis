# Clevis Go

This is a front-end for the [anatol/clevis.go](https://github.com/anatol/clevis.go) module, providing a CLI interface
that attempts to be a Go implementation of [latchset/clevis](https://github.com/latchset/clevis).  
The benefits over the original version being statically-linked binaries for cross platform usage without dependency libraries.

## Usage

```shell
# encrypting
clevis encrypt PIN CONFIG < PLAINTEXT > CIPHERTEXT.jwe

# decrypting
clevis decrypt < CIPHERTEXT.jwe > PLAINTEXT
```

See the [Clevis](https://github.com/latchset/clevis) repo for PINS and CONFIG parameters.
