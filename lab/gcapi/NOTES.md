# Google Calendar API

The Calendar API does return a lot of really useful information, especially for opening events in a
users's browser, like Kairos will aim to do. If you have multiple Google Accounts signed in however,
the links won't work properly.

Each `htmlLink` contains an `eid`. This is the event ID. As long as we have that, and we know which
account ID we want to use (i.e. the account index, starting from 0, ordered by the account you 
added in the browser first to last) then we can create a valid event URL that will actually take
you straight to the event in the calendar.

A URL would look like this:

```
https://calendar.google.com/calendar/b/<ACCOUNT_ID/event?eid=<EVENT_ID>
```

It's not the ideal solution, but it's not the worst either. When setting up Kairos, a user could 
just be told to log in in the same order that they logged in in their browser? There would be a 
settings / credentials file anyway, so a user would be able to change the order if they wanted to.

Otherwise, everything else seems pretty straightforward in reality.
