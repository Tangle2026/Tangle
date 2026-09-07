# Validator subsystem — Alpha

The current alpha records validator fee allocation and validates PoR records through deterministic proof hashes. A production validator set must be implemented inside the consensus layer with:

- validator identity and consensus public keys;
- bonded stake and voting power;
- signed votes/proposals through CometBFT;
- downtime/jailing rules;
- objectively provable slashing conditions;
- deterministic PoR verification;
- validator reward distribution.
