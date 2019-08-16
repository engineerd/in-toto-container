# In-Toto verifications in a container

This is a project that enables running [In-Toto][in-toto] verifications inside a Linux container.

> This project assumes familiarity with the [In-Toto Specification][in-toto-spec].

> Verification is the process by which data and metadata included in the final product is used to ensure its correctness. Verification is performed by the client by checking the supply chain layout and links for correctness, as well as by performing the inspection steps.

## Why run verifications in a container?

There are a couple of reasons for providing a containerized environment for In-Toto verifications:

- dependencies - at a minimum, running an In-Toto verification [requires Python, OpenSSL, git, and pip][deps]. On top of this, individual layouts might require additional dependencies in order to execute the verification (for example, the official demo needs to execute `tar` locally).
- OS independence - we can run verifications from Windows.

[in-toto]: https://github.com/in-toto/in-toto
[in-toto-spec]: https://github.com/in-toto/docs/blob/master/in-toto-spec.md
[deps]: https://github.com/in-toto/in-toto#install-dependencies