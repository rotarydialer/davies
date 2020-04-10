# **davies:** Database Adaptable Variance Idiomatic Evolution Scheduler

## Intended Audience

  - Engineers who regularly manage the creation of scripts to update the
    schema in a MS SQL database.

  - Engineers who want to simplify and/or standardize how other team
    members contribute schema changes to a MS SQL database.

## Purpose

`davies` makes it very simple for engineers to contribute schema changes 
to a MS SQL database, managing the schema evolutions as proper source code. 
Schema changes are deployed as gzipped tarballs named with the corresponding 
git tag.

To apply schema changes to a particular database, pull from git or
download a tarball and use `davies` to automatically apply those 
scripts in chronological order.

`davies` provides simple tools to manage the process of creating and 
applying schema upgrade scripts to databases in all environments.

 - scripts are automatically named with a timestamp assigned at time
   of creation

 - all scripts applied to the database are recorded in
   the table `davies_evo.scripts` - making it simple to
   see what has been applied, and when.

`davies` contains only tools for managing schema evolutions. The idea is
that you create one git repository for each of your databases then use
`davies` to manage the schema evolution of each database.

This project is modeled after 
[Schema Evolution Manager](https://github.com/mbryzek/schema-evolution-manager),
a similar tool for managing schema changes in Postgres. Further, this
project maintains the same goals (see below) and overall philosophy as `sem`,
and can, in fact, be considered more or less a direct port of that tool
for MS SQL. (Though we may not ever reach full feature parity.)

## Project Goals

  - **Absolutely minimal set of dependencies.** We found that anything
    more complex led developers to prefer to manage their own schema
    evolutions. We prefer small sets of scripts that each do one thing
    well.

  - **Committed to true simplicity.** Features that would add complexity
    are simply not added.

  - **Works for ALL applications.** Schema management is a first class
    task now, so any application framework can leverage these
    migration tools.

  - **No rollback.** We have found in practice that rolling back schema
    changes is not 100% reliable. Therefore we inentionally do NOT
    support rollback. This is an often debated aspect of the parent
    project, and although the design itself could be easily extended to 
    support rollback, we have no plans to do so.

In place of rollback, we prefer to keep focus on the criticality of
schema changes, encouraging peer review and lots of smaller evolutions
that themselves are relatively harmless.

This stems from the idea that we believe schema evolutions are
*fundamentally risky*. We believe the best way to manage this risk is
to:

  1. Treat schema evolution changes as normal software releases
  as much as possible.

  2. Manage schema versions as simple tarballs or git branches - 
  artifacts are critical to provide 100% reproducibility. This means 
  the exact same artifacts can be applied in development then QA and 
  finally production environments.

  3. Isolate schema changes as their own deploy. This then
  guarantees that every other application itself can be rolled
  back if needed. In practice, we have seen greater risk when
  applications couple code changes with schema changes.

This last point bears some more detail. By fundamentally deciding to
manage and release schema changes independent of application changes:

  1. Schema changes are required to be incremental. For example, to
  rename a column takes 4 separate, independent production deploys:

    a. add new column
    b. deploy changes in application to use old and new column
    c. remove old column
    d. deploy changes in application to use only new column

  Though at first this may seem more complex, each individual change
  itself is smaller and lower risk.

  2. It is worth repeating that all application deploys can now be
  rolled back. This has been a huge win for our teams.

### Footnote on the name:

This project is named after astrobiologist Paul Davies, inspired 
in part by his paper [The Algorithmic Origins of Life](https://arxiv.org/abs/1207.4803),
which deals with self-organizing systems. While it would be a 
stretch to say we are applying the same principles to databases
in this tool, we do consider this project to be an experiment in 
making databases more adaptable. Davies' thoughts on this resonate
with us.

Oh, and the phrase is totally not a backronym. 😇