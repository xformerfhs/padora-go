# padora (Go)

A program to show how a padding oracle attack works.

[![Go Report Card](https://goreportcard.com/badge/github.com/xformerfhs/padora-go)](https://goreportcard.com/report/github.com/xformerfhs/padora-go)
[![License](https://img.shields.io/github/license/xformerfhs/filesigner)](https://github.com/xformerfhs/filesigner/blob/main/LICENSE)

## Introduction

I have been working in the fascinating field of cryptography for years and regularly hold training courses on the subject for IT architects, programmers, testers, project leads and product owners.

In these training courses, I point out the padding oracle attack, which makes it possible to break a message encrypted with the CBC mode and a suitable padding method without knowing the key or even the encryption method.

The seminar participants are always amazed and ask me how this is supposed to work? 
That started me to build this small demonstration program to show the programmers how such an attack is carried out.

It is very simple and by no means something that can be used in a real-world scenario, but it shows the principles behind this attack.

## Technical background

The [padding oracle attack](https://en.wikipedia.org/wiki/Padding_oracle_attack) on symmetric ciphers was first shown by Serge Vaudenay in [2002](https://www.iacr.org/cryptodb/archive/2002/EUROCRYPT/2850/2850.pdf) and [improved](https://www.usenix.org/legacy/event/sec02/full_papers/black/black_html) by John Black and Hector Urtubia.
It works in the so called [CBC mode](https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#CBC) of [block ciphers](https://en.wikipedia.org/wiki/Block_cipher). 
In addition to the CBC mode, also a highly-structured [padding](https://en.wikipedia.org/wiki/Padding_(cryptography)) is needed.

### Block ciphers

Block ciphers process data in blocks.
In general, the last block is not completely filled with plain text data.
So, to indicate where the plain text ends, the bytes that do not belong to the plain text are filled according to certain methods.
This is called "padding".

### Padding

There is a plentitude of padding methods.
Some are very simple and some are highly structured.
It is these highly structured padding methods that allow a padding oracle.

A general indication of susceptibility to a padding oracle is the existence of a padding error.
If a wrong padding is possible, then the method may be susceptibility to the attack.

### Padding error
The final ingredient is a dedicated error from the recipient if the padding was incorrect.
The name "padding oracle" reflects the fact that the attacker manipulates the encrypted bytes in a certain way and sends them to the receiver.
If the receiver answers with a "wrong padding" error, if the padding is wrong, then the attack is possible.

This resembles questioning an oracle: "Oh, receiver, does the encrypted data I sent you have a valid padding or not?".
The answer (Yes: No padding error, No: Padding error) makes it possible to reconstruct the plain text.

An explicit "wrong padding" return code or message is not even necessary.
It suffices if a padding error can be distinguished from data errors.
E.g. there may be a timing difference.
I.e. if a data error caused by a wrong padding needs less time to show up than a data error caused by invalid data.
This is called a ["timing attack"](https://en.wikipedia.org/wiki/Timing_attack).

### Vulnerable padding methods

However, not all padding methods that show a padding error are susceptible to the attack.

The following padding methods are known to allow a padding oracle attack in combination with CBC mode:

- PKCS#7
- ESP padding ([RFC 4303](https://www.ietf.org/rfc/rfc4303.txt))
- Smartcard padding ([ISO/IEC 7816-4](https://www.iso.org/obp/ui/#iso:std:iso-iec:7816:-4:ed-4:v1:en), [ISO/IEC 9797-1 - Method 2](https://en.wikipedia.org/wiki/ISO/IEC_9797-1)) 

It is interesting to note that there is one padding method, that is **not** vulnerable to a padding oracle attack:
[Arbitrary tail byte padding](https://eprint.iacr.org/2003/098.pdf).

Strangely enough, this secure padding method is not included in any standard or library.
However, it requires a pseudo-random byte, the creation of which is CPU intensive.

### Attack

In CBC mode the previous cipher block is XOR-ed into the current decrypted block.
This means that the attacker can directly change the result of the decryption in a deterministic way!
The attack works by changing the data in the previous block so that the data in the current block is decrypted into a block with a valid padding.
The attacker generates a valid padding where no padding had been before.
From the value that is needed to generate the valid padding the plain text data can be inferred.

The maniplated data will not make sense any more to the receiving application.
But the attacker is only interested in whether there is a padding error, or not.

## Intent

The program that is presented here has been developed to show how such an attack is performed.
It makes it possible to follow the exact steps needed for the attack.
It is intended to further the understanding of the attack in order to make it possible to defend against it.

## Learning

If there is one thing that can be learned from this, it is that encryption must always be combined with authentication.

- **Never** use encryption without authentication.
- **Only** use encryption with authentication.
- Always.
- Even with arbitrary tail byte padding :-).

One may use either ciphers that provide [authenticated encryption](https://en.wikipedia.org/wiki/Authenticated_encryption).
Or, if one needs to use a classic mode like CBC, an explicit authentication with a [Message authentication code](https://en.wikipedia.org/wiki/Message_authentication_code) is mandatory.

If authentication is used the padding method does not matter any more.
One can use whatever method is the most convienent and/or fastest.
Even the vulnerable methods can be used.
As long as any manipulation of the encrypted data can be detected, the vulnerability does not matter any more.

## Contact

Frank Schwab ([Mail](mailto:github.sfdhi@slmails.com "Mail"))

## License

This source code is published under the [Apache License V2](https://www.apache.org/licenses/LICENSE-2.0.txt).
