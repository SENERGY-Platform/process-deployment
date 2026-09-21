# Aspects in a process model are a list

An element of a process model names a set of aspects, not one aspect. The
single-valued fields that predate the list still exist, are deprecated, and
behave as an alias for a list with one element.

## Scope

Applies from `models/go v0.0.0-20260910124809-95949e15e3d4` onwards, to the
aspects this service reads out of a BPMN model and writes into a deployment: the
selection criteria of tasks, message events and conditional events, and the
camunda task payload.

Not about the aspects of a **content variable**. Those live on a device type,
are also spelled `aspect_ids`, and point the other way round: a content variable
enumerates what it carries, while the criteria here demand that all of them be
carried. The two are easy to swap because the field name is the same.

Not about evaluating the criteria either. This service forwards them; the
matching happens in the device-repository and, for conditional events, in the
event pipeline. The only evaluation done here is the service filter described
under "Selection options" below.

## The two spellings and where each is read

| Where | Deprecated | List |
|---|---|---|
| BPMN attribute on a message or conditional event | `senergy:aspect` | `senergy:aspects` |
| camunda task payload | `aspect` | `aspects` |
| `models.ProcessFilterCriteria` | `AspectId` | `AspectIds` |

`senergy:aspects` holds the ids separated by commas — an aspect urn contains
none. `lib/ctrl/deployment/parser/aspects.go` reads both attributes; `hasAspect`
is the gate that decides whether an element counts as an event at all, so it has
to accept the list on its own, or an element that names only the list is not
recognised as an event.

In the task payload both `aspect` and `aspects` are aspect **nodes**, not ids,
because the payload carries the resolved node for the worker.

## Both spellings are passed on as they were selected

Nothing here folds one field into the other. A criteria that names a single
aspect produces `aspect` alone, a criteria that names a list produces `aspects`
alone, and a criteria that carries both produces both.

That is deliberate and it is the reading rule that makes it safe: whoever
evaluates the criteria folds the deprecated field into the list, so a model
written before the lists keeps its results. Folding on the writing side would
put a list into a payload that a reader predating the list cannot see, and drop
the single field that is the only one such a reader looks at.

The consequence for fixtures is worth knowing: because the list fields are
`omitempty`, every payload on a single-aspect path is byte-identical to what it
was before the list existed. A whole-payload fixture only changes where a
process model actually names more than one aspect.

The nodes of `aspects` are resolved in a stable order, sorted by aspect id, so
the same criteria always produces the same payload. Sorting by id also puts the
node that the deprecated single field would carry first, which is how the
platform picks it elsewhere.

## Selection options: several aspects in one criteria are an AND

`serviceMatchesCriteria` in `lib/ctrl/prepare_deployment.go` decides which
services of a selectable device are offered for an element. A criteria naming
more than one aspect requires **one** path option to carry all of them — the
same rule the device-repository applies to a content variable — and each named
aspect covers its own subtree, so a parent aspect matches an option sitting on
one of its children.

Two helpers keep that readable, and they are the two directions of the alias:

- `criteriaAspectIds` folds the criteria's deprecated `AspectId` into its list,
  because a criteria is an input and may carry either spelling. An empty aspect
  is not a filter.
- `pathOptionAspectNodes` reads the answer the other way round: it prefers
  `AspectNodes` and falls back to the single `AspectNode` only when the list is
  empty, because a selection service that predates the list fills the single
  field alone while a current one fills both.

Before the shared model carried `AspectNodes` on its path option, the answer
reached this service narrowed to the alphabetically first node, and a criteria
naming two aspects could not be decided at all.

## The stored deployment names the field differently

`models.ProcessFilterCriteria` carries json tags only, and the deployment is stored
with `bson:",inline"`, so mongo names the field by the driver's own rule rather than
by the wire spelling: `aspectids`, not `aspect_ids`. That is symmetric — the same
rule applies on write and on read, so the list survives a round trip — but it is
invisible from the struct, and a bson tag added to one field and not another would
break it silently. `lib/db/aspects_test.go` asserts the round trip for both
spellings rather than leaving it to be reasoned about.

## Where the aspect list has no producer yet

As of 2026-09-10 no client writes an aspect list into a process model: the
process designer still emits a single `senergy:aspect` and a single `aspect` in
the task payload. Until that changes, an aspect list reaches a deployment only
through a client that posts `filter_criteria.aspect_ids` to `/v3/deployments`
directly. The parsing and stringifying above are in place regardless, so the
model side needs no further change when a producer appears.
