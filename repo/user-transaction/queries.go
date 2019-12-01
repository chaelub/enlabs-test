package user_transaction

const (
	processTransaction = `with transaction as (
    insert into
        transactions(extid,amount,state,status,userid,tms)
        values($1,$2,$3,$4,$5,$6)
    returning
        userid,
        case
          when
             state = 1
          then
             amount*( - 1)
          else
             amount
        end as amount)
update users set account = account+transaction.amount from transaction where id=transaction.userid;`

	cancelLastNTransactions = `with updtr as (update transactions set status = 3 from (
    select
		id,
        row_number() over(order by tms desc) as row,
        userid,
        state,
        amount
    from transactions
    where status = 2
    limit $1
) as st where transactions.id = st.id and mod(st.row,2) = 0
returning st.row, st.userid, st.state, st.amount),
grptr as (
    select
        updtr.userid,
        sum(
               case
                  when
                     updtr.state = 1
                  then
                     updtr.amount*( - 1)
                  else
                     updtr.amount
               end
        ) as amount
    from updtr
    group by updtr.userid)
update
         users
      set
         account = account + grptr.amount
      from grptr
where users.id = grptr.userid;`
)
