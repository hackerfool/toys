-module(binary_tree).

-export([new/0,
		 size/1,
		 insert/2,
		 insert/3,
		 from_list/1,
		 prev/1,
		 mid/1,
		 last/1]).

new() ->
	#{val => nil, left => nil, right => nil}.

size(Tree) ->
	size(Tree, 0).
size(Tree, Size) ->
	Size2 = case maps:get(val, Tree) of
				nil ->
					Size;
				_ ->
					Size + 1
			end,
	Size3 = case maps:get(right, Tree) of
				nil ->
					Size2;
				RTree ->
					size(RTree, Size2)
			end,
	Size4 = case maps:get(left, Tree) of
				nil ->
					Size3;
				LTree ->
					size(LTree, Size3)
			end,
	Size4.

insert(Tree, Value) ->
	CmpFun = fun(A, B) ->
				if
					B == nil ->
						val;
					A < B ->
						left;
					true ->
						right
				end
			 end,
	insert(Tree, Value, CmpFun).

insert(Tree, Value, CmpFun) ->
	Tvalue = maps:get(val, Tree),
	PosFlag = 	CmpFun(Value, Tvalue),
	case maps:get(PosFlag, Tree) of
		nil when PosFlag == val ->
			maps:put(val, Value, Tree);
		nil ->
			NewTree = maps:put(val, Value, new()),
			maps:put(PosFlag, NewTree, Tree);
		Ctree ->
			maps:put(PosFlag, insert(Ctree, Value), Tree)
	end.

from_list([]) ->
	new();
from_list(List) ->
	from_list(new(), List).

from_list(Tree, []) ->
	Tree;
from_list(Tree, [T | L]) ->
	from_list(insert(Tree, T), L).

prev(Tree) ->
	lists:reverse(prev(Tree, [])).
prev(Tree, List) ->
	List2 = case maps:get(val, Tree) of
				nil ->
					List;
				Val ->
					[Val | List]
			end,
	List3 = case maps:get(left, Tree) of
				nil ->
					List2;
				LTree ->
					prev(LTree, List2)
			end,
	List4 = case maps:get(right, Tree) of
				nil ->
					List3;
				RTree ->
					prev(RTree, List3)
			end,
	List4.

mid(Tree) ->
	lists:reverse(mid(Tree, [])).
mid(Tree, List) ->
	List2 = case maps:get(left, Tree) of
				nil ->
					List;
				LTree ->
					mid(LTree, List)
			end,
	List3 = case maps:get(val, Tree) of
				nil ->
					List2;
				Val ->
					[Val | List2]
			end,
	List4 = case maps:get(right, Tree) of
				nil ->
					List3;
				RTree ->
					mid(RTree, List3)
			end,
	List4.

last(Tree) ->
	lists:reverse(last(Tree, [])).
last(Tree, List) ->
	List2 = case maps:get(left, Tree) of
				nil ->
					List;
				LTree ->
					last(LTree, List)
			end,
	List3 = case maps:get(right, Tree) of
				nil ->
					List2;
				RTree ->
					last(RTree, List2)
			end,
	List4 = case maps:get(val, Tree) of
				nil ->
					List3;
				Val ->
					[Val | List3]
			end,
	List4.