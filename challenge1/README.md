SPLICE spec
=======================

| Offset  | Size in bytes (type)    |Description            |
|:-------:|:------------------------|:----------------------|
| 0       | 5 (string)              | SPLICE header         |
| 6       | 8 (int64 big endian)    | Encoded data size     |
| 14      | 33 (string + null bytes)| Hardware version      |
| 46      | 4 (float32)             | BPM                   |
| 50      | encoded data size - 36  | First Track           |

Everything except for the encoded data size is little endian (most likely to throw you off)

Tracks:
* 1 uint32: id
* 1 byte: length of instrument name
* string: instrument name
* [4]uint32 : pattern

The tracks repeat until you've read up to the encoded data size, at which point
you simply stop parsing and return what you have.

A good optimization to make is to read the first parts of the file up until the
encoded data size value, and then read the rest of it and close the file at that point.

This should speed up reading in the cases of decoding files with extra stuff (as seen in
pattern_5.slice).


## Submitting



Create a zip of your Go source code and send the zip file to gochallenge [at] joshsoftware.com
by the 14th of March 2015 (midnight IST, 11:30 AM PDT). No new solutions will be accepted after
that. In the email mention your full name, country of residence, and twitter id (if any). We shall
be publishing on this blog, a list of participant names. If you don't want you name to appear kindly
mention the same in your email. We are accepting anonymous submissions and will evaluate them too
but then these participants are not eligible for the prizes. We will give your zip file to the evaluation
team.

Note: Avoid sharing your code with anyone else; if your solution becomes available to the general public it
might impact evaluation of your submission.

[gochallenge@joshsoftware.com](mailto:gochallenge@johnsoftware.com)

Subject Line: Go Challenge 1 Submission
Name: Full Name
Country of Residence: USA
Twitter: Twitter Handle

