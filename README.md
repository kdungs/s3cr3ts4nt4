# S3cr3ts4nt4

Cryptographically secure (to some extent) gift exchange. Or something along
those lines.

Imagine you want to exchange gifts with your friends in a merry round of secret
Santa. But since you're all remote, you need a way of randomly assigning people
to each other while making sure that

 * only the person who will send a gift to an address gets to see that address
 * people don't get assigned to themselves

This is an attempt to solve this problem with public-key cryptography. The
solution is not perfect, however. It relies on a trusted authority to manage
the gift exchange. The host of the exchange needs to be trusted not to
compromise the system and use the private host key in a malicious way.

Here's how it works.

The host generates a key-pair and sends the public key to all participants. The
participants also generate a key-pair (ideally via a key-derivation algorithm,
so they just need to remember a passphrase and not necessarily store they
keys). They will generate a payload containing their name, address, and public
key. That payload is then signed with their private key, encrypted with the
host's public key, and sent to the host.

Once the host has collected all payloads and verified their signatures, it
assigns each participant a recipient. Technically speaking, it produces a
derangement of the list of participants. For each participant, it uses that
persons public key to encrypt the address of the recipient into a file that can
then be sent to the person.
