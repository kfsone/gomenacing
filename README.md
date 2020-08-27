GoMenacing -- Elite Dangerous Trade Planner
===========================================

GoMenacing is a tool for planning trade activities in the game Elite Dangerous. It tries to give
you answers to the question "I'm here, where can I make money?"

The name is a play on its predecessor, "TradeDangerous". GoM is written in the "Go" language to
give it better performance, and "menacing" is a more active synonym for "dangerous".


Status
======

. Pre-Alpha, Unstable,


Data Import
===========

Instead of making Menacing capable of importing data from many sites, it instead uses a
standard method for declaring a compact, efficient format it wants the data in (protobufs).

This will provide an easy path for implementing "translators" that convert, say, EDDB data
into Menacing. It will also be a huge boon to import and startup times.

https://github.com/kfsone/eddbtrans

