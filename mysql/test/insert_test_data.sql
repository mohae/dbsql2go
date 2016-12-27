--
-- abc
--

TRUNCATE TABLE abc;

INSERT INTO abc (
	code, description, tiny
	, small, medium, ger
	, big, cost, created
)
VALUES (
	"a", "AAA", 1
	, 11, 111, 1111
	, 11111, 11.11, null
);

INSERT INTO abc (
	code, description, tiny
	, small, medium, ger
	, big, cost, created
)
VALUES (
	"b", "BBB", 2
	, 22, 222, 2222
	, 22222, 22.22, null
);

INSERT INTO abc (
	code, description, tiny
	, small, medium, ger
	, big, cost, created
)
VALUES (
	"c", "CCC", 3
	, 33, 333, 3333
	, 33333, 33.33, null
);


INSERT INTO abc (
	code, description, tiny
	, small, medium, ger
	, big, cost, created
)
VALUES (
	"d", "DDD", 4
	, 44, 444, 4444
	, 44444, 44.44, null
);

--
-- abc_nn
--

TRUNCATE TABLE abc_nn;

INSERT INTO abc_nn (
	code, description, tiny
	, small, medium, ger
	, big, cost, created
)
VALUES (
	"a", "AAA", 1
	, 11, 111, 1111
	, 11111, 11.11, null
);

--
-- def
--
TRUNCATE TABLE def;

INSERT INTO def (
	d_year
)
VALUES (
	' 1996'
)
;

INSERT INTO def (
	size,
	a_set
)
VALUES (
	'small',
	'a'
)
;

INSERT INTO def (
	d_date,
	d_datetime,
	d_time,
	d_year,
	size,
	a_set
)
VALUES (
	'2001-12-20',
	'2011-07-04 16:33:22',
	'16:20:00',
	'1999',
	'large',
	'b'
)
;

--
-- def_nn
--
TRUNCATE TABLE def_nn;

INSERT INTO def_nn (
	d_date,
	d_datetime,
	d_time,
	d_year,
	size,
	a_set
)
VALUES (
	'2001-12-20',
	'2011-07-04 16:33:22',
	'16:20:00',
	'1999',
	'large',
	'2'
)
;

--
-- ghi
--
TRUNCATE TABLE ghi;

INSERT INTO ghi (
	tiny_stuff
)
VALUES (
  ' '
)
;

INSERT INTO ghi (
	tiny_stuff
	, stuff
	, med_stuff
	, long_stuff
)
VALUES (
  '0123456789',
  'ABCDEFGHIJKLMNOPQRSTUVWXYZ',
  'abcdefghijklmnopqrstuvwxyz',
  '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'
)
;

--
-- ghi_nn
--
TRUNCATE TABLE ghi_nn;

INSERT INTO ghi_nn (
	tiny_stuff
	, stuff
	, med_stuff
	, long_stuff
)
VALUES (
  '0123456789',
  'ABCDEFGHIJKLMNOPQRSTUVWXYZ',
  'abcdefghijklmnopqrstuvwxyz',
  '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'
)
;

--
-- jkl
--
TRUNCATE TABLE jkl;

INSERT INTO jkl (
	tiny_txt
)
VALUES (
  ' '
)
;

INSERT INTO jkl (
	tiny_txt
	, txt
	, med_txt
	, long_txt
	, bin
	, var_bin
)
VALUES (
  '0123456789',
  'ABCDEFGHIJKLMNOPQRSTUVWXYZ',
  'abcdefghijklmnopqrstuvwxyz',
  '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz',
  ' !#$%&()*+',
  '{|}~'
)
;

--
-- jkl_nn
--
TRUNCATE TABLE jkl_nn;

INSERT INTO jkl_nn (
	tiny_txt
	, txt
	, med_txt
	, long_txt
	, bin
	, var_bin
)
VALUES (
  '0123456789',
  'ABCDEFGHIJKLMNOPQRSTUVWXYZ',
  'abcdefghijklmnopqrstuvwxyz',
  '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz',
  ' !#$%&()*+',
  '{|}~'
)
;
