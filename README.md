# dunakeke.hu

Webpage for DUNAKEKE

# TODOs:

## focus
- [ ] impressum
- [ ] contacts page
- [ ] page for yearly income/stuff...
    - [ ] way to upload pdf files.
- [ ] file/picture selector/search
- [ ] fix donation handling if not returned to site
- [ ] email password reset
    - [ ] normal password reset
- [x] user profiles
    - [ ] passwor change on profile
    - [ ] account delete
- [ ] Post headers

- [ ] Supporter admin update is missing
- [ ] fix per site roles
- [ ] set choose file text
- [ ] import existing pages
- [ ] expost modules
    - [ ] authentication?
    - [ ] config
- [ ] Add missing translations
- [ ] Remove debug messaged
- [ ] Remove TODO-s and FIXME-s

## overall

- [ ] Add donationn refund
- [ ] Add CSRF tokens to cookies..
- [ ] sanitize form lengths
- [ ] move styles to stylesheets instead of template files
- [ ] gitlab diff like editor
- [ ] editable menubar
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
- [ ] implement home page..
- [ ] implement admin pages..
- [ ] Check all links are working
- [ ] Use release contexts
- [ ] Make everything fancier
- [ ] fix editor translations
- [ ] use statistics
- [ ] make comments work

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
- [x] check if saving per-user site statistics hurts GDPR: they do not
- [x] post tags
    - [x] should show up on front page
    - [x] tags should prevent public listing of posts
- [x] add simplepay logo
- [x] add supporter and their logos...
- [x] show page Titles properly
- [x] remove TODO() contexts
- [x] Implement post editing..
- [x] email services
- [x] Donations fail to update on site, if the user closes the simplepay site.
- [x] Show simplepay ID
- [x] File upload error... -> Generate files folder if does not exists...
