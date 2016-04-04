##Rest Interface

Resource		| GET			|	POST | PUT | DELETE
-------		| -------	   |	-----| ----| ------
/customer|	(200)List of customers|	404| 404 | 404
/customer/{token}/users	|(200) List of user objects| 404| 404| 404
/customer/{token}/users/{userid} |	(200) single user object|404|(200) activate	| (200) deactivate
/customer/{token}/users/{email} |	(200) single user object | 404 |	(200) activate	| (200) deactivate
/customer/{token}/groups |	(200) List of groups |404 | 404 |404
/customer/{token}/groups/{groupid} |	(200) List of members |(200) add user |404	 |(200) deactivate user