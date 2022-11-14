-- name: FindLastActivityByAccount :one
  select activity.*
    from account_activity as activity 
   where activity.account_uuid = $1  
order by activity.date_time desc
   limit 1;

-- name: FindAccountStatuses :one
  select activity.activity_type = 'ENABLED' and activity.activity_type <> 'INVALIDATED' as IS_ACTIVE, 
         activity.activity_type = 'DISABLED' as IS_DISABLED,
         activity.activity_type = 'INVALIDATED' as IS_INVALIDATED
    from account_activity as activity 
   where activity.account_uuid = '283191b8-3bc0-4b9a-9be9-e44635eb0438'  
order by activity.date_time desc
  limit 1;
