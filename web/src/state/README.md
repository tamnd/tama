# state

Client state is hand-rolled React context, no store library.

M1 carries exactly one piece of cross-screen state: the session user that `GET /api/me` returns.
A context provider with two setters covers that in about sixty lines, keeps the web dependency list at react, react-dom, and the router, and leaves nothing to learn, version, or upgrade.
The server is the source of truth anyway; the cookie is the session, so the client never persists anything.

If a later milestone grows real cross-screen state (course progress, an in-flight lesson, settings), revisit zustand then.
Swapping a context read for a store hook is a small, local change, so deciding late costs nothing now.
