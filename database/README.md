# database
Not sure how this will be used yet.  Putting instructions here for later use.

## Setup
Instructions for installing postgres and creating a database.

```
# Get the pq library
go get github.com/lib/pq

# Install postgres
sudo -i

  yum install postgresql-server postgresql-contrib
  postgresql-setup initdb

  # need to enable password authentication
  vi /var/lib/pgsql/data/pg_hba.conf
  # At the bottom, instead of ending in ident, the following lines should end in md5, like so:
    host    all             all             127.0.0.1/32            md5
    host    all             all             ::1/128                 md5

  # Start the service and create a symlink
  systemctl start postgresql
  systemctl enable postgresql

  # No longer need to be root.
  exit

# Change to the postgres user:
sudo -i -u postgres

  # Open a postgres prompt:
  psql
    # Opens a postgres prompt (confirm it works).
    # Quit the prompt
    \q

  # Create new role:
  createuser --interactive
    dev (or whatever your preferred non-root linux username is)
    y (superuser)

  # Create a database with your new username
  createdb dev (or whatever the username is)

  # No longer need to be the postgres user.
  exit

# As your user you just named a new database after (probably dev)
psql

  # To connect to a different database, you'd type psql -d (database name)
  # Get some user and connection info:
  \conninfo

  # Update your password
  ALTER USER dev PASSWORD 'dev';   # If dev is the username you're working with and you want a bad password.

  # Create a table
  CREATE TABLE widgets (
      widgetid serial PRIMARY KEY,
      widgetname varchar(50) NOT NULL,
      widgetsomething varchar(50) NOT NULL
  );

  # Add some groups and their keys (This is just an example.  More to come...)
  INSERT INTO widgets (widgetname, widgetsomething) VALUES ('Super Widget', '1234567890');
  INSERT INTO widgets (widgetname, widgetsomething) VALUES ('Mighty Fine Widget', 'abcdefghij');

  # Confirm the groups and keys were added.
  SELECT * FROM widgets;

  # Quit the psql prompt
  \q
```
