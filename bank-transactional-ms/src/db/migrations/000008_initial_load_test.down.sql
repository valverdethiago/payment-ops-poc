
delete 
  from account_activity 
 where account_uuid in (
          select account_uuid 
            from account 
           where bank_uuid = (select bank_uuid from bank where code = '001') 
             and account_number in ('00000001', '00000002', '00000003')
     );
     
delete 
  from account
 where bank_uuid = (select bank_uuid from bank where code = '001') 
   and account_number in ('00000001', '00000002', '00000003');
   
delete 
  from configuration
 where bank_uuid = (select bank_uuid from bank where code = '001') ;

delete
  from bank 
 where code = '001';
