insert into bank (code, name, country_code)
values ('001', 'Banco do Brasil', 'BR');

insert into configuration (bank_uuid, kafka_input_topic_name)
values ( (select bank_uuid from bank where code = '001'), 'TRIO_SYNC_REQUEST_INPUT' );


insert into account (account_uuid, account_number, account_type, bank_uuid)
values ('9b2272fe-53ad-4d5d-bfaf-ce2f1cf27ccf', '00000001', 'CHECKING', (select bank_uuid from bank where code = '001'));

insert into account_activity (account_uuid, activity_type, date_time)
values ('9b2272fe-53ad-4d5d-bfaf-ce2f1cf27ccf', 'CREATED', TO_TIMESTAMP('2022-05-26 9:30:20', 'YYYY-MM-DD HH:MI:SS') );

insert into account_activity (account_uuid, activity_type, date_time)
values ('9b2272fe-53ad-4d5d-bfaf-ce2f1cf27ccf', 'ENABLED', TO_TIMESTAMP('2022-05-26 9:30:30', 'YYYY-MM-DD HH:MI:SS') );


insert into account (account_uuid, account_number, account_type, bank_uuid)
values ('86c36433-f7ef-4a03-8dda-8690d7168521', '00000002', 'CHECKING', (select bank_uuid from bank where code = '001'));

insert into account_activity (account_uuid, activity_type, date_time)
values ('86c36433-f7ef-4a03-8dda-8690d7168521', 'CREATED', TO_TIMESTAMP('2022-05-26 9:30:20', 'YYYY-MM-DD HH:MI:SS') );

insert into account_activity (account_uuid, activity_type, date_time)
values ('86c36433-f7ef-4a03-8dda-8690d7168521', 'ENABLED', TO_TIMESTAMP('2022-05-26 9:30:30', 'YYYY-MM-DD HH:MI:SS') );

insert into account_activity (account_uuid, activity_type, date_time)
values ('86c36433-f7ef-4a03-8dda-8690d7168521', 'DISABLED', TO_TIMESTAMP('2022-05-26 9:30:40', 'YYYY-MM-DD HH:MI:SS') );


insert into account (account_uuid, account_number, account_type, bank_uuid)
values ('465cff98-24d4-4f88-baa0-944ff8c97bb0', '00000003', 'CHECKING', (select bank_uuid from bank where code = '001'));

insert into account_activity (account_uuid, activity_type, date_time)
values ( '465cff98-24d4-4f88-baa0-944ff8c97bb0', 'CREATED', TO_TIMESTAMP('2022-05-26 9:30:20', 'YYYY-MM-DD HH:MI:SS') );

insert into account_activity (account_uuid, activity_type, date_time)
values ( '465cff98-24d4-4f88-baa0-944ff8c97bb0', 'ENABLED', TO_TIMESTAMP('2022-05-26 9:30:30', 'YYYY-MM-DD HH:MI:SS') );

insert into account_activity (account_uuid, activity_type, date_time)
values ( '465cff98-24d4-4f88-baa0-944ff8c97bb0', 'INVALIDATED', TO_TIMESTAMP('2022-05-26 9:30:40', 'YYYY-MM-DD HH:MI:SS') );