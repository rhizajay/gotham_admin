#gotham_admin

Goal of this project is to build some basic administration tools for gotham using Go.

##Subprojects

Feature		| Package			|	Description
-------		| -------	|	-----------
gtmadmin	| [gtadmin](https://)		| A basic command line tool. 
gtmadsrv	| [web](https://)	|Web service
gtmadweb	| [web](https://)	|Web client for gotham_administration
gtmslack	| [slack](https://)	|Slackbot


##Libraries

Feature		| Package			|	Description
-------		| -------	|	-----------
Users	| [Users](https://)		| List User; Search by email, Search by Id
Groups	| [Groups](https://)	| List Groups; List Members; add user; Remove user
Tasks	|[Tasks](https://)		| List tasks; Change Task Status
Reports|[Reports](https://)	| List Reports; Generate Stats

##Why Admin tools?
Starting off to scratch my own itch of managing and mainting user accounts.  This should make it easier for other users to make changes.  Proof of Concept creating services around existing infrastructure to improve performance.
