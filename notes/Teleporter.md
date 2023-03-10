# Teleporter

Lower registers control things like sound and light
Eight register controls teleportation

Two destinations:
* used when eigth register is at its minimum energy level (default operation)
* requires specific energy level otehrwise not recommended

Need to extract confirmation algorithm, reimplement, and optimize it.
Set the eighth register, activate teleporter, bypass confirmation.

by default, register 7 (8) is always 0 and is read from twice

first after reading `teleporter`
`eq Reading value 114 from register 4`

`call 6027` seems important, right before that, `set r0 4` `set r1 1`, right after that `eq r1 r0 6`

using a hammer to just set r0 to 6 with an arbitrary r7 gives bad code (e.g. `uXWjUDMpMpyd` for 1000), code must be derived from r7. BUT you end up on the beach, so I could potentially pursue code 8 instead
```
r0, r1 = 4, 1
func myFunc() {
  if r0 != 0 {
    if r1 != 0 {
      stack.push(r0)
      r1 = r1 + 32767
      myFunc()
      r1 = r0
      r0 = stack.pop()
      r0 = r0 + 32767
      myFunc()
      return
    }
    r0 = r0 + 32767
    r1 = r7
    myFunc()
    return
  }
  r0 = r1 + 1
  return
}
```

with a different value
```
A strange, electronic voice is projected into your mind:

  "Unusual setting detected!  Starting confirmation process!  Estimated time to completion: 1 billion years."
```
look for `(")` in logs to find algo
reg 1 inc by 1
`pop \(3\)\n32769\nret \(18\)\nadd` search in memoryDump

```
The cover of this book subtly swirls with colors.  It is titled "A Brief Introduction to Interdimensional Physics".  It reads:

Recent advances in interdimensional physics have produced fascinating
predictions about the fundamentals of our universe!  For example,
interdimensional physics seems to predict that the universe is, at its root, a
purely mathematical construct, and that all events are caused by the
interactions between eight pockets of energy called "registers".
Furthermore, it seems that while the lower registers primarily control mundane
things like sound and light, the highest register (the so-called "eighth
register") is used to control interdimensional events such as teleportation.

A hypothetical such teleportation device would need to have have exactly two
destinations.  One destination would be used when the eighth register is at its
minimum energy level - this would be the default operation assuming the user
has no way to control the eighth register.  In this situation, the teleporter
should send the user to a preconfigured safe location as a default.

The second destination, however, is predicted to require a very specific
energy level in the eighth register.  The teleporter must take great care to
confirm that this energy level is exactly correct before teleporting its user!
If it is even slightly off, the user would (probably) arrive at the correct
location, but would briefly experience anomalies in the fabric of reality
itself - this is, of course, not recommended.  Any teleporter would need to test
the energy level in the eighth register and abort teleportation if it is not
exactly correct.

This required precision implies that the confirmation mechanism would be very
computationally expensive.  While this would likely not be an issue for large-
scale teleporters, a hypothetical hand-held teleporter would take billions of
years to compute the result and confirm that the eighth register is correct.

If you find yourself trapped in an alternate dimension with nothing but a
hand-held teleporter, you will need to extract the confirmation algorithm,
reimplement it on more powerful hardware, and optimize it.  This should, at the
very least, allow you to determine the value of the eighth register which would
have been accepted by the teleporter's confirmation mechanism.

Then, set the eighth register to this value, activate the teleporter, and
bypass the confirmation mechanism.  If the eighth register is set correctly, no
anomalies should be experienced, but beware - if it is set incorrectly, the
now-bypassed confirmation mechanism will not protect you!

Of course, since teleportation is impossible, this is all totally ridiculous.
```