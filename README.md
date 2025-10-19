# dunakeke.hu

Webpage for DUNAKEKE

# TODOs:

## focus
- [ ] fix per site roles
- [ ] set choose file text
- [ ] check if saving per-user site statistics hurts GDPR
- [ ] import existing pages
- [ ] email password reset
- [ ] expost modules
    - [ ] authentication?
    - [ ] config

## overall

- [ ] improve routing with getting the routes from somewhere...
- [ ] set admin page access for roles
- [ ] make crud options <templates>
- [ ] alternatives input validation to match root link format
- [ ] Sanitize inputs!
- [ ] add delete user option
- [x] Donation backend
    - [ ] Check Donation backend
- [ ] Handle link redirecting...
- [ ] Proper messaging between states
- [ ] Add artifacts and other directories to config
- [ ] Image gallery and selection
- [ ] Everything else
- [ ] User role add
- [ ] Make all cruds editable from the UI
- [x] Link otpay correctly...
    - [ ] look into applepay??
    - [ ] Final check
- [ ] Implement post editing..
- [ ] implement home page..
- [ ] implement admin pages..
- [ ] Check all links are working
- [ ] Use release contexts
- [ ] Make everything fancier
- [ ] fix editor translations
- [ ] use statistics

- [ ] Add statistics
    - [ ] Log statistics from every site opening
        - [ ] From
        - [ ] Else
    - [ ] Show it on admins page

## done

- [x] proper post view
- [x] proper post photo editing
- [x] post route translation
- [x] hidden post handling
- [x] user role handling
- [x] expost modules
    - [x] rendering
    - [x] logger
    - [x] sessions
- [x] Add statistics
    - [x] Log statistics from every site opening
        - [x] Logging
- [x] Fix page not found translation
- [x] Link otpay correctly...
- [x] Donation backend
- [x] random login failure bug and password hash deletion fix
    - [x] user logic dto didn't contain the password hash, so it avoided to save it back on update
