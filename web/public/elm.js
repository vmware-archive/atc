
(function() {
'use strict';

function F2(fun)
{
  function wrapper(a) { return function(b) { return fun(a,b); }; }
  wrapper.arity = 2;
  wrapper.func = fun;
  return wrapper;
}

function F3(fun)
{
  function wrapper(a) {
    return function(b) { return function(c) { return fun(a, b, c); }; };
  }
  wrapper.arity = 3;
  wrapper.func = fun;
  return wrapper;
}

function F4(fun)
{
  function wrapper(a) { return function(b) { return function(c) {
    return function(d) { return fun(a, b, c, d); }; }; };
  }
  wrapper.arity = 4;
  wrapper.func = fun;
  return wrapper;
}

function F5(fun)
{
  function wrapper(a) { return function(b) { return function(c) {
    return function(d) { return function(e) { return fun(a, b, c, d, e); }; }; }; };
  }
  wrapper.arity = 5;
  wrapper.func = fun;
  return wrapper;
}

function F6(fun)
{
  function wrapper(a) { return function(b) { return function(c) {
    return function(d) { return function(e) { return function(f) {
    return fun(a, b, c, d, e, f); }; }; }; }; };
  }
  wrapper.arity = 6;
  wrapper.func = fun;
  return wrapper;
}

function F7(fun)
{
  function wrapper(a) { return function(b) { return function(c) {
    return function(d) { return function(e) { return function(f) {
    return function(g) { return fun(a, b, c, d, e, f, g); }; }; }; }; }; };
  }
  wrapper.arity = 7;
  wrapper.func = fun;
  return wrapper;
}

function F8(fun)
{
  function wrapper(a) { return function(b) { return function(c) {
    return function(d) { return function(e) { return function(f) {
    return function(g) { return function(h) {
    return fun(a, b, c, d, e, f, g, h); }; }; }; }; }; }; };
  }
  wrapper.arity = 8;
  wrapper.func = fun;
  return wrapper;
}

function F9(fun)
{
  function wrapper(a) { return function(b) { return function(c) {
    return function(d) { return function(e) { return function(f) {
    return function(g) { return function(h) { return function(i) {
    return fun(a, b, c, d, e, f, g, h, i); }; }; }; }; }; }; }; };
  }
  wrapper.arity = 9;
  wrapper.func = fun;
  return wrapper;
}

function A2(fun, a, b)
{
  return fun.arity === 2
    ? fun.func(a, b)
    : fun(a)(b);
}
function A3(fun, a, b, c)
{
  return fun.arity === 3
    ? fun.func(a, b, c)
    : fun(a)(b)(c);
}
function A4(fun, a, b, c, d)
{
  return fun.arity === 4
    ? fun.func(a, b, c, d)
    : fun(a)(b)(c)(d);
}
function A5(fun, a, b, c, d, e)
{
  return fun.arity === 5
    ? fun.func(a, b, c, d, e)
    : fun(a)(b)(c)(d)(e);
}
function A6(fun, a, b, c, d, e, f)
{
  return fun.arity === 6
    ? fun.func(a, b, c, d, e, f)
    : fun(a)(b)(c)(d)(e)(f);
}
function A7(fun, a, b, c, d, e, f, g)
{
  return fun.arity === 7
    ? fun.func(a, b, c, d, e, f, g)
    : fun(a)(b)(c)(d)(e)(f)(g);
}
function A8(fun, a, b, c, d, e, f, g, h)
{
  return fun.arity === 8
    ? fun.func(a, b, c, d, e, f, g, h)
    : fun(a)(b)(c)(d)(e)(f)(g)(h);
}
function A9(fun, a, b, c, d, e, f, g, h, i)
{
  return fun.arity === 9
    ? fun.func(a, b, c, d, e, f, g, h, i)
    : fun(a)(b)(c)(d)(e)(f)(g)(h)(i);
}

//import Native.Utils //

var _elm_lang$core$Native_Basics = function() {

function div(a, b)
{
	return (a / b) | 0;
}
function rem(a, b)
{
	return a % b;
}
function mod(a, b)
{
	if (b === 0)
	{
		throw new Error('Cannot perform mod 0. Division by zero error.');
	}
	var r = a % b;
	var m = a === 0 ? 0 : (b > 0 ? (a >= 0 ? r : r + b) : -mod(-a, -b));

	return m === b ? 0 : m;
}
function logBase(base, n)
{
	return Math.log(n) / Math.log(base);
}
function negate(n)
{
	return -n;
}
function abs(n)
{
	return n < 0 ? -n : n;
}

function min(a, b)
{
	return _elm_lang$core$Native_Utils.cmp(a, b) < 0 ? a : b;
}
function max(a, b)
{
	return _elm_lang$core$Native_Utils.cmp(a, b) > 0 ? a : b;
}
function clamp(lo, hi, n)
{
	return _elm_lang$core$Native_Utils.cmp(n, lo) < 0
		? lo
		: _elm_lang$core$Native_Utils.cmp(n, hi) > 0
			? hi
			: n;
}

var ord = ['LT', 'EQ', 'GT'];

function compare(x, y)
{
	return { ctor: ord[_elm_lang$core$Native_Utils.cmp(x, y) + 1] };
}

function xor(a, b)
{
	return a !== b;
}
function not(b)
{
	return !b;
}
function isInfinite(n)
{
	return n === Infinity || n === -Infinity;
}

function truncate(n)
{
	return n | 0;
}

function degrees(d)
{
	return d * Math.PI / 180;
}
function turns(t)
{
	return 2 * Math.PI * t;
}
function fromPolar(point)
{
	var r = point._0;
	var t = point._1;
	return _elm_lang$core$Native_Utils.Tuple2(r * Math.cos(t), r * Math.sin(t));
}
function toPolar(point)
{
	var x = point._0;
	var y = point._1;
	return _elm_lang$core$Native_Utils.Tuple2(Math.sqrt(x * x + y * y), Math.atan2(y, x));
}

return {
	div: F2(div),
	rem: F2(rem),
	mod: F2(mod),

	pi: Math.PI,
	e: Math.E,
	cos: Math.cos,
	sin: Math.sin,
	tan: Math.tan,
	acos: Math.acos,
	asin: Math.asin,
	atan: Math.atan,
	atan2: F2(Math.atan2),

	degrees: degrees,
	turns: turns,
	fromPolar: fromPolar,
	toPolar: toPolar,

	sqrt: Math.sqrt,
	logBase: F2(logBase),
	negate: negate,
	abs: abs,
	min: F2(min),
	max: F2(max),
	clamp: F3(clamp),
	compare: F2(compare),

	xor: F2(xor),
	not: not,

	truncate: truncate,
	ceiling: Math.ceil,
	floor: Math.floor,
	round: Math.round,
	toFloat: function(x) { return x; },
	isNaN: isNaN,
	isInfinite: isInfinite
};

}();
//import //

var _elm_lang$core$Native_Utils = function() {

// COMPARISONS

function eq(x, y)
{
	var stack = [];
	var isEqual = eqHelp(x, y, 0, stack);
	var pair;
	while (isEqual && (pair = stack.pop()))
	{
		isEqual = eqHelp(pair.x, pair.y, 0, stack);
	}
	return isEqual;
}


function eqHelp(x, y, depth, stack)
{
	if (depth > 100)
	{
		stack.push({ x: x, y: y });
		return true;
	}

	if (x === y)
	{
		return true;
	}

	if (typeof x !== 'object')
	{
		if (typeof x === 'function')
		{
			throw new Error(
				'Trying to use `(==)` on functions. There is no way to know if functions are "the same" in the Elm sense.'
				+ ' Read more about this at http://package.elm-lang.org/packages/elm-lang/core/latest/Basics#=='
				+ ' which describes why it is this way and what the better version will look like.'
			);
		}
		return false;
	}

	if (x === null || y === null)
	{
		return false
	}

	if (x instanceof Date)
	{
		return x.getTime() === y.getTime();
	}

	if (!('ctor' in x))
	{
		for (var key in x)
		{
			if (!eqHelp(x[key], y[key], depth + 1, stack))
			{
				return false;
			}
		}
		return true;
	}

	// convert Dicts and Sets to lists
	if (x.ctor === 'RBNode_elm_builtin' || x.ctor === 'RBEmpty_elm_builtin')
	{
		x = _elm_lang$core$Dict$toList(x);
		y = _elm_lang$core$Dict$toList(y);
	}
	if (x.ctor === 'Set_elm_builtin')
	{
		x = _elm_lang$core$Set$toList(x);
		y = _elm_lang$core$Set$toList(y);
	}

	// check if lists are equal without recursion
	if (x.ctor === '::')
	{
		var a = x;
		var b = y;
		while (a.ctor === '::' && b.ctor === '::')
		{
			if (!eqHelp(a._0, b._0, depth + 1, stack))
			{
				return false;
			}
			a = a._1;
			b = b._1;
		}
		return a.ctor === b.ctor;
	}

	// check if Arrays are equal
	if (x.ctor === '_Array')
	{
		var xs = _elm_lang$core$Native_Array.toJSArray(x);
		var ys = _elm_lang$core$Native_Array.toJSArray(y);
		if (xs.length !== ys.length)
		{
			return false;
		}
		for (var i = 0; i < xs.length; i++)
		{
			if (!eqHelp(xs[i], ys[i], depth + 1, stack))
			{
				return false;
			}
		}
		return true;
	}

	if (!eqHelp(x.ctor, y.ctor, depth + 1, stack))
	{
		return false;
	}

	for (var key in x)
	{
		if (!eqHelp(x[key], y[key], depth + 1, stack))
		{
			return false;
		}
	}
	return true;
}

// Code in Generate/JavaScript.hs, Basics.js, and List.js depends on
// the particular integer values assigned to LT, EQ, and GT.

var LT = -1, EQ = 0, GT = 1;

function cmp(x, y)
{
	if (typeof x !== 'object')
	{
		return x === y ? EQ : x < y ? LT : GT;
	}

	if (x instanceof String)
	{
		var a = x.valueOf();
		var b = y.valueOf();
		return a === b ? EQ : a < b ? LT : GT;
	}

	if (x.ctor === '::' || x.ctor === '[]')
	{
		while (x.ctor === '::' && y.ctor === '::')
		{
			var ord = cmp(x._0, y._0);
			if (ord !== EQ)
			{
				return ord;
			}
			x = x._1;
			y = y._1;
		}
		return x.ctor === y.ctor ? EQ : x.ctor === '[]' ? LT : GT;
	}

	if (x.ctor.slice(0, 6) === '_Tuple')
	{
		var ord;
		var n = x.ctor.slice(6) - 0;
		var err = 'cannot compare tuples with more than 6 elements.';
		if (n === 0) return EQ;
		if (n >= 1) { ord = cmp(x._0, y._0); if (ord !== EQ) return ord;
		if (n >= 2) { ord = cmp(x._1, y._1); if (ord !== EQ) return ord;
		if (n >= 3) { ord = cmp(x._2, y._2); if (ord !== EQ) return ord;
		if (n >= 4) { ord = cmp(x._3, y._3); if (ord !== EQ) return ord;
		if (n >= 5) { ord = cmp(x._4, y._4); if (ord !== EQ) return ord;
		if (n >= 6) { ord = cmp(x._5, y._5); if (ord !== EQ) return ord;
		if (n >= 7) throw new Error('Comparison error: ' + err); } } } } } }
		return EQ;
	}

	throw new Error(
		'Comparison error: comparison is only defined on ints, '
		+ 'floats, times, chars, strings, lists of comparable values, '
		+ 'and tuples of comparable values.'
	);
}


// COMMON VALUES

var Tuple0 = {
	ctor: '_Tuple0'
};

function Tuple2(x, y)
{
	return {
		ctor: '_Tuple2',
		_0: x,
		_1: y
	};
}

function chr(c)
{
	return new String(c);
}


// GUID

var count = 0;
function guid(_)
{
	return count++;
}


// RECORDS

function update(oldRecord, updatedFields)
{
	var newRecord = {};
	for (var key in oldRecord)
	{
		var value = (key in updatedFields) ? updatedFields[key] : oldRecord[key];
		newRecord[key] = value;
	}
	return newRecord;
}


//// LIST STUFF ////

var Nil = { ctor: '[]' };

function Cons(hd, tl)
{
	return {
		ctor: '::',
		_0: hd,
		_1: tl
	};
}

function append(xs, ys)
{
	// append Strings
	if (typeof xs === 'string')
	{
		return xs + ys;
	}

	// append Lists
	if (xs.ctor === '[]')
	{
		return ys;
	}
	var root = Cons(xs._0, Nil);
	var curr = root;
	xs = xs._1;
	while (xs.ctor !== '[]')
	{
		curr._1 = Cons(xs._0, Nil);
		xs = xs._1;
		curr = curr._1;
	}
	curr._1 = ys;
	return root;
}


// CRASHES

function crash(moduleName, region)
{
	return function(message) {
		throw new Error(
			'Ran into a `Debug.crash` in module `' + moduleName + '` ' + regionToString(region) + '\n'
			+ 'The message provided by the code author is:\n\n    '
			+ message
		);
	};
}

function crashCase(moduleName, region, value)
{
	return function(message) {
		throw new Error(
			'Ran into a `Debug.crash` in module `' + moduleName + '`\n\n'
			+ 'This was caused by the `case` expression ' + regionToString(region) + '.\n'
			+ 'One of the branches ended with a crash and the following value got through:\n\n    ' + toString(value) + '\n\n'
			+ 'The message provided by the code author is:\n\n    '
			+ message
		);
	};
}

function regionToString(region)
{
	if (region.start.line == region.end.line)
	{
		return 'on line ' + region.start.line;
	}
	return 'between lines ' + region.start.line + ' and ' + region.end.line;
}


// TO STRING

function toString(v)
{
	var type = typeof v;
	if (type === 'function')
	{
		var name = v.func ? v.func.name : v.name;
		return '<function' + (name === '' ? '' : ':') + name + '>';
	}

	if (type === 'boolean')
	{
		return v ? 'True' : 'False';
	}

	if (type === 'number')
	{
		return v + '';
	}

	if (v instanceof String)
	{
		return '\'' + addSlashes(v, true) + '\'';
	}

	if (type === 'string')
	{
		return '"' + addSlashes(v, false) + '"';
	}

	if (v === null)
	{
		return 'null';
	}

	if (type === 'object' && 'ctor' in v)
	{
		var ctorStarter = v.ctor.substring(0, 5);

		if (ctorStarter === '_Tupl')
		{
			var output = [];
			for (var k in v)
			{
				if (k === 'ctor') continue;
				output.push(toString(v[k]));
			}
			return '(' + output.join(',') + ')';
		}

		if (ctorStarter === '_Task')
		{
			return '<task>'
		}

		if (v.ctor === '_Array')
		{
			var list = _elm_lang$core$Array$toList(v);
			return 'Array.fromList ' + toString(list);
		}

		if (v.ctor === '<decoder>')
		{
			return '<decoder>';
		}

		if (v.ctor === '_Process')
		{
			return '<process:' + v.id + '>';
		}

		if (v.ctor === '::')
		{
			var output = '[' + toString(v._0);
			v = v._1;
			while (v.ctor === '::')
			{
				output += ',' + toString(v._0);
				v = v._1;
			}
			return output + ']';
		}

		if (v.ctor === '[]')
		{
			return '[]';
		}

		if (v.ctor === 'Set_elm_builtin')
		{
			return 'Set.fromList ' + toString(_elm_lang$core$Set$toList(v));
		}

		if (v.ctor === 'RBNode_elm_builtin' || v.ctor === 'RBEmpty_elm_builtin')
		{
			return 'Dict.fromList ' + toString(_elm_lang$core$Dict$toList(v));
		}

		var output = '';
		for (var i in v)
		{
			if (i === 'ctor') continue;
			var str = toString(v[i]);
			var c0 = str[0];
			var parenless = c0 === '{' || c0 === '(' || c0 === '<' || c0 === '"' || str.indexOf(' ') < 0;
			output += ' ' + (parenless ? str : '(' + str + ')');
		}
		return v.ctor + output;
	}

	if (type === 'object')
	{
		if (v instanceof Date)
		{
			return '<' + v.toString() + '>';
		}

		if (v.elm_web_socket)
		{
			return '<websocket>';
		}

		var output = [];
		for (var k in v)
		{
			output.push(k + ' = ' + toString(v[k]));
		}
		if (output.length === 0)
		{
			return '{}';
		}
		return '{ ' + output.join(', ') + ' }';
	}

	return '<internal structure>';
}

function addSlashes(str, isChar)
{
	var s = str.replace(/\\/g, '\\\\')
			  .replace(/\n/g, '\\n')
			  .replace(/\t/g, '\\t')
			  .replace(/\r/g, '\\r')
			  .replace(/\v/g, '\\v')
			  .replace(/\0/g, '\\0');
	if (isChar)
	{
		return s.replace(/\'/g, '\\\'');
	}
	else
	{
		return s.replace(/\"/g, '\\"');
	}
}


return {
	eq: eq,
	cmp: cmp,
	Tuple0: Tuple0,
	Tuple2: Tuple2,
	chr: chr,
	update: update,
	guid: guid,

	append: F2(append),

	crash: crash,
	crashCase: crashCase,

	toString: toString
};

}();
var _elm_lang$core$Basics$uncurry = F2(
	function (f, _p0) {
		var _p1 = _p0;
		return A2(f, _p1._0, _p1._1);
	});
var _elm_lang$core$Basics$curry = F3(
	function (f, a, b) {
		return f(
			{ctor: '_Tuple2', _0: a, _1: b});
	});
var _elm_lang$core$Basics$flip = F3(
	function (f, b, a) {
		return A2(f, a, b);
	});
var _elm_lang$core$Basics$snd = function (_p2) {
	var _p3 = _p2;
	return _p3._1;
};
var _elm_lang$core$Basics$fst = function (_p4) {
	var _p5 = _p4;
	return _p5._0;
};
var _elm_lang$core$Basics$always = F2(
	function (a, _p6) {
		return a;
	});
var _elm_lang$core$Basics$identity = function (x) {
	return x;
};
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['<|'] = F2(
	function (f, x) {
		return f(x);
	});
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['|>'] = F2(
	function (x, f) {
		return f(x);
	});
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['>>'] = F3(
	function (f, g, x) {
		return g(
			f(x));
	});
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['<<'] = F3(
	function (g, f, x) {
		return g(
			f(x));
	});
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['++'] = _elm_lang$core$Native_Utils.append;
var _elm_lang$core$Basics$toString = _elm_lang$core$Native_Utils.toString;
var _elm_lang$core$Basics$isInfinite = _elm_lang$core$Native_Basics.isInfinite;
var _elm_lang$core$Basics$isNaN = _elm_lang$core$Native_Basics.isNaN;
var _elm_lang$core$Basics$toFloat = _elm_lang$core$Native_Basics.toFloat;
var _elm_lang$core$Basics$ceiling = _elm_lang$core$Native_Basics.ceiling;
var _elm_lang$core$Basics$floor = _elm_lang$core$Native_Basics.floor;
var _elm_lang$core$Basics$truncate = _elm_lang$core$Native_Basics.truncate;
var _elm_lang$core$Basics$round = _elm_lang$core$Native_Basics.round;
var _elm_lang$core$Basics$not = _elm_lang$core$Native_Basics.not;
var _elm_lang$core$Basics$xor = _elm_lang$core$Native_Basics.xor;
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['||'] = _elm_lang$core$Native_Basics.or;
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['&&'] = _elm_lang$core$Native_Basics.and;
var _elm_lang$core$Basics$max = _elm_lang$core$Native_Basics.max;
var _elm_lang$core$Basics$min = _elm_lang$core$Native_Basics.min;
var _elm_lang$core$Basics$compare = _elm_lang$core$Native_Basics.compare;
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['>='] = _elm_lang$core$Native_Basics.ge;
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['<='] = _elm_lang$core$Native_Basics.le;
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['>'] = _elm_lang$core$Native_Basics.gt;
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['<'] = _elm_lang$core$Native_Basics.lt;
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['/='] = _elm_lang$core$Native_Basics.neq;
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['=='] = _elm_lang$core$Native_Basics.eq;
var _elm_lang$core$Basics$e = _elm_lang$core$Native_Basics.e;
var _elm_lang$core$Basics$pi = _elm_lang$core$Native_Basics.pi;
var _elm_lang$core$Basics$clamp = _elm_lang$core$Native_Basics.clamp;
var _elm_lang$core$Basics$logBase = _elm_lang$core$Native_Basics.logBase;
var _elm_lang$core$Basics$abs = _elm_lang$core$Native_Basics.abs;
var _elm_lang$core$Basics$negate = _elm_lang$core$Native_Basics.negate;
var _elm_lang$core$Basics$sqrt = _elm_lang$core$Native_Basics.sqrt;
var _elm_lang$core$Basics$atan2 = _elm_lang$core$Native_Basics.atan2;
var _elm_lang$core$Basics$atan = _elm_lang$core$Native_Basics.atan;
var _elm_lang$core$Basics$asin = _elm_lang$core$Native_Basics.asin;
var _elm_lang$core$Basics$acos = _elm_lang$core$Native_Basics.acos;
var _elm_lang$core$Basics$tan = _elm_lang$core$Native_Basics.tan;
var _elm_lang$core$Basics$sin = _elm_lang$core$Native_Basics.sin;
var _elm_lang$core$Basics$cos = _elm_lang$core$Native_Basics.cos;
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['^'] = _elm_lang$core$Native_Basics.exp;
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['%'] = _elm_lang$core$Native_Basics.mod;
var _elm_lang$core$Basics$rem = _elm_lang$core$Native_Basics.rem;
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['//'] = _elm_lang$core$Native_Basics.div;
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['/'] = _elm_lang$core$Native_Basics.floatDiv;
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['*'] = _elm_lang$core$Native_Basics.mul;
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['-'] = _elm_lang$core$Native_Basics.sub;
var _elm_lang$core$Basics_ops = _elm_lang$core$Basics_ops || {};
_elm_lang$core$Basics_ops['+'] = _elm_lang$core$Native_Basics.add;
var _elm_lang$core$Basics$toPolar = _elm_lang$core$Native_Basics.toPolar;
var _elm_lang$core$Basics$fromPolar = _elm_lang$core$Native_Basics.fromPolar;
var _elm_lang$core$Basics$turns = _elm_lang$core$Native_Basics.turns;
var _elm_lang$core$Basics$degrees = _elm_lang$core$Native_Basics.degrees;
var _elm_lang$core$Basics$radians = function (t) {
	return t;
};
var _elm_lang$core$Basics$GT = {ctor: 'GT'};
var _elm_lang$core$Basics$EQ = {ctor: 'EQ'};
var _elm_lang$core$Basics$LT = {ctor: 'LT'};
var _elm_lang$core$Basics$Never = function (a) {
	return {ctor: 'Never', _0: a};
};

//import Native.Utils //

var _elm_lang$core$Native_Debug = function() {

function log(tag, value)
{
	var msg = tag + ': ' + _elm_lang$core$Native_Utils.toString(value);
	var process = process || {};
	if (process.stdout)
	{
		process.stdout.write(msg);
	}
	else
	{
		console.log(msg);
	}
	return value;
}

function crash(message)
{
	throw new Error(message);
}

return {
	crash: crash,
	log: F2(log)
};

}();
var _elm_lang$core$Debug$crash = _elm_lang$core$Native_Debug.crash;
var _elm_lang$core$Debug$log = _elm_lang$core$Native_Debug.log;

var _elm_lang$core$Maybe$withDefault = F2(
	function ($default, maybe) {
		var _p0 = maybe;
		if (_p0.ctor === 'Just') {
			return _p0._0;
		} else {
			return $default;
		}
	});
var _elm_lang$core$Maybe$Nothing = {ctor: 'Nothing'};
var _elm_lang$core$Maybe$oneOf = function (maybes) {
	oneOf:
	while (true) {
		var _p1 = maybes;
		if (_p1.ctor === '[]') {
			return _elm_lang$core$Maybe$Nothing;
		} else {
			var _p3 = _p1._0;
			var _p2 = _p3;
			if (_p2.ctor === 'Nothing') {
				var _v3 = _p1._1;
				maybes = _v3;
				continue oneOf;
			} else {
				return _p3;
			}
		}
	}
};
var _elm_lang$core$Maybe$andThen = F2(
	function (maybeValue, callback) {
		var _p4 = maybeValue;
		if (_p4.ctor === 'Just') {
			return callback(_p4._0);
		} else {
			return _elm_lang$core$Maybe$Nothing;
		}
	});
var _elm_lang$core$Maybe$Just = function (a) {
	return {ctor: 'Just', _0: a};
};
var _elm_lang$core$Maybe$map = F2(
	function (f, maybe) {
		var _p5 = maybe;
		if (_p5.ctor === 'Just') {
			return _elm_lang$core$Maybe$Just(
				f(_p5._0));
		} else {
			return _elm_lang$core$Maybe$Nothing;
		}
	});
var _elm_lang$core$Maybe$map2 = F3(
	function (func, ma, mb) {
		var _p6 = {ctor: '_Tuple2', _0: ma, _1: mb};
		if (((_p6.ctor === '_Tuple2') && (_p6._0.ctor === 'Just')) && (_p6._1.ctor === 'Just')) {
			return _elm_lang$core$Maybe$Just(
				A2(func, _p6._0._0, _p6._1._0));
		} else {
			return _elm_lang$core$Maybe$Nothing;
		}
	});
var _elm_lang$core$Maybe$map3 = F4(
	function (func, ma, mb, mc) {
		var _p7 = {ctor: '_Tuple3', _0: ma, _1: mb, _2: mc};
		if ((((_p7.ctor === '_Tuple3') && (_p7._0.ctor === 'Just')) && (_p7._1.ctor === 'Just')) && (_p7._2.ctor === 'Just')) {
			return _elm_lang$core$Maybe$Just(
				A3(func, _p7._0._0, _p7._1._0, _p7._2._0));
		} else {
			return _elm_lang$core$Maybe$Nothing;
		}
	});
var _elm_lang$core$Maybe$map4 = F5(
	function (func, ma, mb, mc, md) {
		var _p8 = {ctor: '_Tuple4', _0: ma, _1: mb, _2: mc, _3: md};
		if (((((_p8.ctor === '_Tuple4') && (_p8._0.ctor === 'Just')) && (_p8._1.ctor === 'Just')) && (_p8._2.ctor === 'Just')) && (_p8._3.ctor === 'Just')) {
			return _elm_lang$core$Maybe$Just(
				A4(func, _p8._0._0, _p8._1._0, _p8._2._0, _p8._3._0));
		} else {
			return _elm_lang$core$Maybe$Nothing;
		}
	});
var _elm_lang$core$Maybe$map5 = F6(
	function (func, ma, mb, mc, md, me) {
		var _p9 = {ctor: '_Tuple5', _0: ma, _1: mb, _2: mc, _3: md, _4: me};
		if ((((((_p9.ctor === '_Tuple5') && (_p9._0.ctor === 'Just')) && (_p9._1.ctor === 'Just')) && (_p9._2.ctor === 'Just')) && (_p9._3.ctor === 'Just')) && (_p9._4.ctor === 'Just')) {
			return _elm_lang$core$Maybe$Just(
				A5(func, _p9._0._0, _p9._1._0, _p9._2._0, _p9._3._0, _p9._4._0));
		} else {
			return _elm_lang$core$Maybe$Nothing;
		}
	});

//import Native.Utils //

var _elm_lang$core$Native_List = function() {

var Nil = { ctor: '[]' };

function Cons(hd, tl)
{
	return { ctor: '::', _0: hd, _1: tl };
}

function fromArray(arr)
{
	var out = Nil;
	for (var i = arr.length; i--; )
	{
		out = Cons(arr[i], out);
	}
	return out;
}

function toArray(xs)
{
	var out = [];
	while (xs.ctor !== '[]')
	{
		out.push(xs._0);
		xs = xs._1;
	}
	return out;
}


function range(lo, hi)
{
	var list = Nil;
	if (lo <= hi)
	{
		do
		{
			list = Cons(hi, list);
		}
		while (hi-- > lo);
	}
	return list;
}

function foldr(f, b, xs)
{
	var arr = toArray(xs);
	var acc = b;
	for (var i = arr.length; i--; )
	{
		acc = A2(f, arr[i], acc);
	}
	return acc;
}

function map2(f, xs, ys)
{
	var arr = [];
	while (xs.ctor !== '[]' && ys.ctor !== '[]')
	{
		arr.push(A2(f, xs._0, ys._0));
		xs = xs._1;
		ys = ys._1;
	}
	return fromArray(arr);
}

function map3(f, xs, ys, zs)
{
	var arr = [];
	while (xs.ctor !== '[]' && ys.ctor !== '[]' && zs.ctor !== '[]')
	{
		arr.push(A3(f, xs._0, ys._0, zs._0));
		xs = xs._1;
		ys = ys._1;
		zs = zs._1;
	}
	return fromArray(arr);
}

function map4(f, ws, xs, ys, zs)
{
	var arr = [];
	while (   ws.ctor !== '[]'
		   && xs.ctor !== '[]'
		   && ys.ctor !== '[]'
		   && zs.ctor !== '[]')
	{
		arr.push(A4(f, ws._0, xs._0, ys._0, zs._0));
		ws = ws._1;
		xs = xs._1;
		ys = ys._1;
		zs = zs._1;
	}
	return fromArray(arr);
}

function map5(f, vs, ws, xs, ys, zs)
{
	var arr = [];
	while (   vs.ctor !== '[]'
		   && ws.ctor !== '[]'
		   && xs.ctor !== '[]'
		   && ys.ctor !== '[]'
		   && zs.ctor !== '[]')
	{
		arr.push(A5(f, vs._0, ws._0, xs._0, ys._0, zs._0));
		vs = vs._1;
		ws = ws._1;
		xs = xs._1;
		ys = ys._1;
		zs = zs._1;
	}
	return fromArray(arr);
}

function sortBy(f, xs)
{
	return fromArray(toArray(xs).sort(function(a, b) {
		return _elm_lang$core$Native_Utils.cmp(f(a), f(b));
	}));
}

function sortWith(f, xs)
{
	return fromArray(toArray(xs).sort(function(a, b) {
		var ord = f(a)(b).ctor;
		return ord === 'EQ' ? 0 : ord === 'LT' ? -1 : 1;
	}));
}

return {
	Nil: Nil,
	Cons: Cons,
	cons: F2(Cons),
	toArray: toArray,
	fromArray: fromArray,
	range: range,

	foldr: F3(foldr),

	map2: F3(map2),
	map3: F4(map3),
	map4: F5(map4),
	map5: F6(map5),
	sortBy: F2(sortBy),
	sortWith: F2(sortWith)
};

}();
var _elm_lang$core$List$sortWith = _elm_lang$core$Native_List.sortWith;
var _elm_lang$core$List$sortBy = _elm_lang$core$Native_List.sortBy;
var _elm_lang$core$List$sort = function (xs) {
	return A2(_elm_lang$core$List$sortBy, _elm_lang$core$Basics$identity, xs);
};
var _elm_lang$core$List$drop = F2(
	function (n, list) {
		drop:
		while (true) {
			if (_elm_lang$core$Native_Utils.cmp(n, 0) < 1) {
				return list;
			} else {
				var _p0 = list;
				if (_p0.ctor === '[]') {
					return list;
				} else {
					var _v1 = n - 1,
						_v2 = _p0._1;
					n = _v1;
					list = _v2;
					continue drop;
				}
			}
		}
	});
var _elm_lang$core$List$map5 = _elm_lang$core$Native_List.map5;
var _elm_lang$core$List$map4 = _elm_lang$core$Native_List.map4;
var _elm_lang$core$List$map3 = _elm_lang$core$Native_List.map3;
var _elm_lang$core$List$map2 = _elm_lang$core$Native_List.map2;
var _elm_lang$core$List$any = F2(
	function (isOkay, list) {
		any:
		while (true) {
			var _p1 = list;
			if (_p1.ctor === '[]') {
				return false;
			} else {
				if (isOkay(_p1._0)) {
					return true;
				} else {
					var _v4 = isOkay,
						_v5 = _p1._1;
					isOkay = _v4;
					list = _v5;
					continue any;
				}
			}
		}
	});
var _elm_lang$core$List$all = F2(
	function (isOkay, list) {
		return _elm_lang$core$Basics$not(
			A2(
				_elm_lang$core$List$any,
				function (_p2) {
					return _elm_lang$core$Basics$not(
						isOkay(_p2));
				},
				list));
	});
var _elm_lang$core$List$foldr = _elm_lang$core$Native_List.foldr;
var _elm_lang$core$List$foldl = F3(
	function (func, acc, list) {
		foldl:
		while (true) {
			var _p3 = list;
			if (_p3.ctor === '[]') {
				return acc;
			} else {
				var _v7 = func,
					_v8 = A2(func, _p3._0, acc),
					_v9 = _p3._1;
				func = _v7;
				acc = _v8;
				list = _v9;
				continue foldl;
			}
		}
	});
var _elm_lang$core$List$length = function (xs) {
	return A3(
		_elm_lang$core$List$foldl,
		F2(
			function (_p4, i) {
				return i + 1;
			}),
		0,
		xs);
};
var _elm_lang$core$List$sum = function (numbers) {
	return A3(
		_elm_lang$core$List$foldl,
		F2(
			function (x, y) {
				return x + y;
			}),
		0,
		numbers);
};
var _elm_lang$core$List$product = function (numbers) {
	return A3(
		_elm_lang$core$List$foldl,
		F2(
			function (x, y) {
				return x * y;
			}),
		1,
		numbers);
};
var _elm_lang$core$List$maximum = function (list) {
	var _p5 = list;
	if (_p5.ctor === '::') {
		return _elm_lang$core$Maybe$Just(
			A3(_elm_lang$core$List$foldl, _elm_lang$core$Basics$max, _p5._0, _p5._1));
	} else {
		return _elm_lang$core$Maybe$Nothing;
	}
};
var _elm_lang$core$List$minimum = function (list) {
	var _p6 = list;
	if (_p6.ctor === '::') {
		return _elm_lang$core$Maybe$Just(
			A3(_elm_lang$core$List$foldl, _elm_lang$core$Basics$min, _p6._0, _p6._1));
	} else {
		return _elm_lang$core$Maybe$Nothing;
	}
};
var _elm_lang$core$List$indexedMap = F2(
	function (f, xs) {
		return A3(
			_elm_lang$core$List$map2,
			f,
			_elm_lang$core$Native_List.range(
				0,
				_elm_lang$core$List$length(xs) - 1),
			xs);
	});
var _elm_lang$core$List$member = F2(
	function (x, xs) {
		return A2(
			_elm_lang$core$List$any,
			function (a) {
				return _elm_lang$core$Native_Utils.eq(a, x);
			},
			xs);
	});
var _elm_lang$core$List$isEmpty = function (xs) {
	var _p7 = xs;
	if (_p7.ctor === '[]') {
		return true;
	} else {
		return false;
	}
};
var _elm_lang$core$List$tail = function (list) {
	var _p8 = list;
	if (_p8.ctor === '::') {
		return _elm_lang$core$Maybe$Just(_p8._1);
	} else {
		return _elm_lang$core$Maybe$Nothing;
	}
};
var _elm_lang$core$List$head = function (list) {
	var _p9 = list;
	if (_p9.ctor === '::') {
		return _elm_lang$core$Maybe$Just(_p9._0);
	} else {
		return _elm_lang$core$Maybe$Nothing;
	}
};
var _elm_lang$core$List_ops = _elm_lang$core$List_ops || {};
_elm_lang$core$List_ops['::'] = _elm_lang$core$Native_List.cons;
var _elm_lang$core$List$map = F2(
	function (f, xs) {
		return A3(
			_elm_lang$core$List$foldr,
			F2(
				function (x, acc) {
					return A2(
						_elm_lang$core$List_ops['::'],
						f(x),
						acc);
				}),
			_elm_lang$core$Native_List.fromArray(
				[]),
			xs);
	});
var _elm_lang$core$List$filter = F2(
	function (pred, xs) {
		var conditionalCons = F2(
			function (x, xs$) {
				return pred(x) ? A2(_elm_lang$core$List_ops['::'], x, xs$) : xs$;
			});
		return A3(
			_elm_lang$core$List$foldr,
			conditionalCons,
			_elm_lang$core$Native_List.fromArray(
				[]),
			xs);
	});
var _elm_lang$core$List$maybeCons = F3(
	function (f, mx, xs) {
		var _p10 = f(mx);
		if (_p10.ctor === 'Just') {
			return A2(_elm_lang$core$List_ops['::'], _p10._0, xs);
		} else {
			return xs;
		}
	});
var _elm_lang$core$List$filterMap = F2(
	function (f, xs) {
		return A3(
			_elm_lang$core$List$foldr,
			_elm_lang$core$List$maybeCons(f),
			_elm_lang$core$Native_List.fromArray(
				[]),
			xs);
	});
var _elm_lang$core$List$reverse = function (list) {
	return A3(
		_elm_lang$core$List$foldl,
		F2(
			function (x, y) {
				return A2(_elm_lang$core$List_ops['::'], x, y);
			}),
		_elm_lang$core$Native_List.fromArray(
			[]),
		list);
};
var _elm_lang$core$List$scanl = F3(
	function (f, b, xs) {
		var scan1 = F2(
			function (x, accAcc) {
				var _p11 = accAcc;
				if (_p11.ctor === '::') {
					return A2(
						_elm_lang$core$List_ops['::'],
						A2(f, x, _p11._0),
						accAcc);
				} else {
					return _elm_lang$core$Native_List.fromArray(
						[]);
				}
			});
		return _elm_lang$core$List$reverse(
			A3(
				_elm_lang$core$List$foldl,
				scan1,
				_elm_lang$core$Native_List.fromArray(
					[b]),
				xs));
	});
var _elm_lang$core$List$append = F2(
	function (xs, ys) {
		var _p12 = ys;
		if (_p12.ctor === '[]') {
			return xs;
		} else {
			return A3(
				_elm_lang$core$List$foldr,
				F2(
					function (x, y) {
						return A2(_elm_lang$core$List_ops['::'], x, y);
					}),
				ys,
				xs);
		}
	});
var _elm_lang$core$List$concat = function (lists) {
	return A3(
		_elm_lang$core$List$foldr,
		_elm_lang$core$List$append,
		_elm_lang$core$Native_List.fromArray(
			[]),
		lists);
};
var _elm_lang$core$List$concatMap = F2(
	function (f, list) {
		return _elm_lang$core$List$concat(
			A2(_elm_lang$core$List$map, f, list));
	});
var _elm_lang$core$List$partition = F2(
	function (pred, list) {
		var step = F2(
			function (x, _p13) {
				var _p14 = _p13;
				var _p16 = _p14._0;
				var _p15 = _p14._1;
				return pred(x) ? {
					ctor: '_Tuple2',
					_0: A2(_elm_lang$core$List_ops['::'], x, _p16),
					_1: _p15
				} : {
					ctor: '_Tuple2',
					_0: _p16,
					_1: A2(_elm_lang$core$List_ops['::'], x, _p15)
				};
			});
		return A3(
			_elm_lang$core$List$foldr,
			step,
			{
				ctor: '_Tuple2',
				_0: _elm_lang$core$Native_List.fromArray(
					[]),
				_1: _elm_lang$core$Native_List.fromArray(
					[])
			},
			list);
	});
var _elm_lang$core$List$unzip = function (pairs) {
	var step = F2(
		function (_p18, _p17) {
			var _p19 = _p18;
			var _p20 = _p17;
			return {
				ctor: '_Tuple2',
				_0: A2(_elm_lang$core$List_ops['::'], _p19._0, _p20._0),
				_1: A2(_elm_lang$core$List_ops['::'], _p19._1, _p20._1)
			};
		});
	return A3(
		_elm_lang$core$List$foldr,
		step,
		{
			ctor: '_Tuple2',
			_0: _elm_lang$core$Native_List.fromArray(
				[]),
			_1: _elm_lang$core$Native_List.fromArray(
				[])
		},
		pairs);
};
var _elm_lang$core$List$intersperse = F2(
	function (sep, xs) {
		var _p21 = xs;
		if (_p21.ctor === '[]') {
			return _elm_lang$core$Native_List.fromArray(
				[]);
		} else {
			var step = F2(
				function (x, rest) {
					return A2(
						_elm_lang$core$List_ops['::'],
						sep,
						A2(_elm_lang$core$List_ops['::'], x, rest));
				});
			var spersed = A3(
				_elm_lang$core$List$foldr,
				step,
				_elm_lang$core$Native_List.fromArray(
					[]),
				_p21._1);
			return A2(_elm_lang$core$List_ops['::'], _p21._0, spersed);
		}
	});
var _elm_lang$core$List$takeReverse = F3(
	function (n, list, taken) {
		takeReverse:
		while (true) {
			if (_elm_lang$core$Native_Utils.cmp(n, 0) < 1) {
				return taken;
			} else {
				var _p22 = list;
				if (_p22.ctor === '[]') {
					return taken;
				} else {
					var _v23 = n - 1,
						_v24 = _p22._1,
						_v25 = A2(_elm_lang$core$List_ops['::'], _p22._0, taken);
					n = _v23;
					list = _v24;
					taken = _v25;
					continue takeReverse;
				}
			}
		}
	});
var _elm_lang$core$List$takeTailRec = F2(
	function (n, list) {
		return _elm_lang$core$List$reverse(
			A3(
				_elm_lang$core$List$takeReverse,
				n,
				list,
				_elm_lang$core$Native_List.fromArray(
					[])));
	});
var _elm_lang$core$List$takeFast = F3(
	function (ctr, n, list) {
		if (_elm_lang$core$Native_Utils.cmp(n, 0) < 1) {
			return _elm_lang$core$Native_List.fromArray(
				[]);
		} else {
			var _p23 = {ctor: '_Tuple2', _0: n, _1: list};
			_v26_5:
			do {
				_v26_1:
				do {
					if (_p23.ctor === '_Tuple2') {
						if (_p23._1.ctor === '[]') {
							return list;
						} else {
							if (_p23._1._1.ctor === '::') {
								switch (_p23._0) {
									case 1:
										break _v26_1;
									case 2:
										return _elm_lang$core$Native_List.fromArray(
											[_p23._1._0, _p23._1._1._0]);
									case 3:
										if (_p23._1._1._1.ctor === '::') {
											return _elm_lang$core$Native_List.fromArray(
												[_p23._1._0, _p23._1._1._0, _p23._1._1._1._0]);
										} else {
											break _v26_5;
										}
									default:
										if ((_p23._1._1._1.ctor === '::') && (_p23._1._1._1._1.ctor === '::')) {
											var _p28 = _p23._1._1._1._0;
											var _p27 = _p23._1._1._0;
											var _p26 = _p23._1._0;
											var _p25 = _p23._1._1._1._1._0;
											var _p24 = _p23._1._1._1._1._1;
											return (_elm_lang$core$Native_Utils.cmp(ctr, 1000) > 0) ? A2(
												_elm_lang$core$List_ops['::'],
												_p26,
												A2(
													_elm_lang$core$List_ops['::'],
													_p27,
													A2(
														_elm_lang$core$List_ops['::'],
														_p28,
														A2(
															_elm_lang$core$List_ops['::'],
															_p25,
															A2(_elm_lang$core$List$takeTailRec, n - 4, _p24))))) : A2(
												_elm_lang$core$List_ops['::'],
												_p26,
												A2(
													_elm_lang$core$List_ops['::'],
													_p27,
													A2(
														_elm_lang$core$List_ops['::'],
														_p28,
														A2(
															_elm_lang$core$List_ops['::'],
															_p25,
															A3(_elm_lang$core$List$takeFast, ctr + 1, n - 4, _p24)))));
										} else {
											break _v26_5;
										}
								}
							} else {
								if (_p23._0 === 1) {
									break _v26_1;
								} else {
									break _v26_5;
								}
							}
						}
					} else {
						break _v26_5;
					}
				} while(false);
				return _elm_lang$core$Native_List.fromArray(
					[_p23._1._0]);
			} while(false);
			return list;
		}
	});
var _elm_lang$core$List$take = F2(
	function (n, list) {
		return A3(_elm_lang$core$List$takeFast, 0, n, list);
	});
var _elm_lang$core$List$repeatHelp = F3(
	function (result, n, value) {
		repeatHelp:
		while (true) {
			if (_elm_lang$core$Native_Utils.cmp(n, 0) < 1) {
				return result;
			} else {
				var _v27 = A2(_elm_lang$core$List_ops['::'], value, result),
					_v28 = n - 1,
					_v29 = value;
				result = _v27;
				n = _v28;
				value = _v29;
				continue repeatHelp;
			}
		}
	});
var _elm_lang$core$List$repeat = F2(
	function (n, value) {
		return A3(
			_elm_lang$core$List$repeatHelp,
			_elm_lang$core$Native_List.fromArray(
				[]),
			n,
			value);
	});

var _elm_lang$core$Result$toMaybe = function (result) {
	var _p0 = result;
	if (_p0.ctor === 'Ok') {
		return _elm_lang$core$Maybe$Just(_p0._0);
	} else {
		return _elm_lang$core$Maybe$Nothing;
	}
};
var _elm_lang$core$Result$withDefault = F2(
	function (def, result) {
		var _p1 = result;
		if (_p1.ctor === 'Ok') {
			return _p1._0;
		} else {
			return def;
		}
	});
var _elm_lang$core$Result$Err = function (a) {
	return {ctor: 'Err', _0: a};
};
var _elm_lang$core$Result$andThen = F2(
	function (result, callback) {
		var _p2 = result;
		if (_p2.ctor === 'Ok') {
			return callback(_p2._0);
		} else {
			return _elm_lang$core$Result$Err(_p2._0);
		}
	});
var _elm_lang$core$Result$Ok = function (a) {
	return {ctor: 'Ok', _0: a};
};
var _elm_lang$core$Result$map = F2(
	function (func, ra) {
		var _p3 = ra;
		if (_p3.ctor === 'Ok') {
			return _elm_lang$core$Result$Ok(
				func(_p3._0));
		} else {
			return _elm_lang$core$Result$Err(_p3._0);
		}
	});
var _elm_lang$core$Result$map2 = F3(
	function (func, ra, rb) {
		var _p4 = {ctor: '_Tuple2', _0: ra, _1: rb};
		if (_p4._0.ctor === 'Ok') {
			if (_p4._1.ctor === 'Ok') {
				return _elm_lang$core$Result$Ok(
					A2(func, _p4._0._0, _p4._1._0));
			} else {
				return _elm_lang$core$Result$Err(_p4._1._0);
			}
		} else {
			return _elm_lang$core$Result$Err(_p4._0._0);
		}
	});
var _elm_lang$core$Result$map3 = F4(
	function (func, ra, rb, rc) {
		var _p5 = {ctor: '_Tuple3', _0: ra, _1: rb, _2: rc};
		if (_p5._0.ctor === 'Ok') {
			if (_p5._1.ctor === 'Ok') {
				if (_p5._2.ctor === 'Ok') {
					return _elm_lang$core$Result$Ok(
						A3(func, _p5._0._0, _p5._1._0, _p5._2._0));
				} else {
					return _elm_lang$core$Result$Err(_p5._2._0);
				}
			} else {
				return _elm_lang$core$Result$Err(_p5._1._0);
			}
		} else {
			return _elm_lang$core$Result$Err(_p5._0._0);
		}
	});
var _elm_lang$core$Result$map4 = F5(
	function (func, ra, rb, rc, rd) {
		var _p6 = {ctor: '_Tuple4', _0: ra, _1: rb, _2: rc, _3: rd};
		if (_p6._0.ctor === 'Ok') {
			if (_p6._1.ctor === 'Ok') {
				if (_p6._2.ctor === 'Ok') {
					if (_p6._3.ctor === 'Ok') {
						return _elm_lang$core$Result$Ok(
							A4(func, _p6._0._0, _p6._1._0, _p6._2._0, _p6._3._0));
					} else {
						return _elm_lang$core$Result$Err(_p6._3._0);
					}
				} else {
					return _elm_lang$core$Result$Err(_p6._2._0);
				}
			} else {
				return _elm_lang$core$Result$Err(_p6._1._0);
			}
		} else {
			return _elm_lang$core$Result$Err(_p6._0._0);
		}
	});
var _elm_lang$core$Result$map5 = F6(
	function (func, ra, rb, rc, rd, re) {
		var _p7 = {ctor: '_Tuple5', _0: ra, _1: rb, _2: rc, _3: rd, _4: re};
		if (_p7._0.ctor === 'Ok') {
			if (_p7._1.ctor === 'Ok') {
				if (_p7._2.ctor === 'Ok') {
					if (_p7._3.ctor === 'Ok') {
						if (_p7._4.ctor === 'Ok') {
							return _elm_lang$core$Result$Ok(
								A5(func, _p7._0._0, _p7._1._0, _p7._2._0, _p7._3._0, _p7._4._0));
						} else {
							return _elm_lang$core$Result$Err(_p7._4._0);
						}
					} else {
						return _elm_lang$core$Result$Err(_p7._3._0);
					}
				} else {
					return _elm_lang$core$Result$Err(_p7._2._0);
				}
			} else {
				return _elm_lang$core$Result$Err(_p7._1._0);
			}
		} else {
			return _elm_lang$core$Result$Err(_p7._0._0);
		}
	});
var _elm_lang$core$Result$formatError = F2(
	function (f, result) {
		var _p8 = result;
		if (_p8.ctor === 'Ok') {
			return _elm_lang$core$Result$Ok(_p8._0);
		} else {
			return _elm_lang$core$Result$Err(
				f(_p8._0));
		}
	});
var _elm_lang$core$Result$fromMaybe = F2(
	function (err, maybe) {
		var _p9 = maybe;
		if (_p9.ctor === 'Just') {
			return _elm_lang$core$Result$Ok(_p9._0);
		} else {
			return _elm_lang$core$Result$Err(err);
		}
	});

//import //

var _elm_lang$core$Native_Platform = function() {


// PROGRAMS

function addPublicModule(object, name, main)
{
	var init = main ? makeEmbed(name, main) : mainIsUndefined(name);

	object['worker'] = function worker(flags)
	{
		return init(undefined, flags, false);
	}

	object['embed'] = function embed(domNode, flags)
	{
		return init(domNode, flags, true);
	}

	object['fullscreen'] = function fullscreen(flags)
	{
		return init(document.body, flags, true);
	};
}


// PROGRAM FAIL

function mainIsUndefined(name)
{
	return function(domNode)
	{
		var message = 'Cannot initialize module `' + name +
			'` because it has no `main` value!\nWhat should I show on screen?';
		domNode.innerHTML = errorHtml(message);
		throw new Error(message);
	};
}

function errorHtml(message)
{
	return '<div style="padding-left:1em;">'
		+ '<h2 style="font-weight:normal;"><b>Oops!</b> Something went wrong when starting your Elm program.</h2>'
		+ '<pre style="padding-left:1em;">' + message + '</pre>'
		+ '</div>';
}


// PROGRAM SUCCESS

function makeEmbed(moduleName, main)
{
	return function embed(rootDomNode, flags, withRenderer)
	{
		try
		{
			var program = mainToProgram(moduleName, main);
			if (!withRenderer)
			{
				program.renderer = dummyRenderer;
			}
			return makeEmbedHelp(moduleName, program, rootDomNode, flags);
		}
		catch (e)
		{
			rootDomNode.innerHTML = errorHtml(e.message);
			throw e;
		}
	};
}

function dummyRenderer()
{
	return { update: function() {} };
}


// MAIN TO PROGRAM

function mainToProgram(moduleName, wrappedMain)
{
	var main = wrappedMain.main;

	if (typeof main.init === 'undefined')
	{
		var emptyBag = batch(_elm_lang$core$Native_List.Nil);
		var noChange = _elm_lang$core$Native_Utils.Tuple2(
			_elm_lang$core$Native_Utils.Tuple0,
			emptyBag
		);

		return _elm_lang$virtual_dom$VirtualDom$programWithFlags({
			init: function() { return noChange; },
			view: function() { return main; },
			update: F2(function() { return noChange; }),
			subscriptions: function () { return emptyBag; }
		});
	}

	var flags = wrappedMain.flags;
	var init = flags
		? initWithFlags(moduleName, main.init, flags)
		: initWithoutFlags(moduleName, main.init);

	return _elm_lang$virtual_dom$VirtualDom$programWithFlags({
		init: init,
		view: main.view,
		update: main.update,
		subscriptions: main.subscriptions,
	});
}

function initWithoutFlags(moduleName, realInit)
{
	return function init(flags)
	{
		if (typeof flags !== 'undefined')
		{
			throw new Error(
				'You are giving module `' + moduleName + '` an argument in JavaScript.\n'
				+ 'This module does not take arguments though! You probably need to change the\n'
				+ 'initialization code to something like `Elm.' + moduleName + '.fullscreen()`'
			);
		}
		return realInit();
	};
}

function initWithFlags(moduleName, realInit, flagDecoder)
{
	return function init(flags)
	{
		var result = A2(_elm_lang$core$Native_Json.run, flagDecoder, flags);
		if (result.ctor === 'Err')
		{
			throw new Error(
				'You are trying to initialize module `' + moduleName + '` with an unexpected argument.\n'
				+ 'When trying to convert it to a usable Elm value, I run into this problem:\n\n'
				+ result._0
			);
		}
		return realInit(result._0);
	};
}


// SETUP RUNTIME SYSTEM

function makeEmbedHelp(moduleName, program, rootDomNode, flags)
{
	var init = program.init;
	var update = program.update;
	var subscriptions = program.subscriptions;
	var view = program.view;
	var makeRenderer = program.renderer;

	// ambient state
	var managers = {};
	var renderer;

	// init and update state in main process
	var initApp = _elm_lang$core$Native_Scheduler.nativeBinding(function(callback) {
		var results = init(flags);
		var model = results._0;
		renderer = makeRenderer(rootDomNode, enqueue, view(model));
		var cmds = results._1;
		var subs = subscriptions(model);
		dispatchEffects(managers, cmds, subs);
		callback(_elm_lang$core$Native_Scheduler.succeed(model));
	});

	function onMessage(msg, model)
	{
		return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback) {
			var results = A2(update, msg, model);
			model = results._0;
			renderer.update(view(model));
			var cmds = results._1;
			var subs = subscriptions(model);
			dispatchEffects(managers, cmds, subs);
			callback(_elm_lang$core$Native_Scheduler.succeed(model));
		});
	}

	var mainProcess = spawnLoop(initApp, onMessage);

	function enqueue(msg)
	{
		_elm_lang$core$Native_Scheduler.rawSend(mainProcess, msg);
	}

	var ports = setupEffects(managers, enqueue);

	return ports ? { ports: ports } : {};
}


// EFFECT MANAGERS

var effectManagers = {};

function setupEffects(managers, callback)
{
	var ports;

	// setup all necessary effect managers
	for (var key in effectManagers)
	{
		var manager = effectManagers[key];

		if (manager.isForeign)
		{
			ports = ports || {};
			ports[key] = manager.tag === 'cmd'
				? setupOutgoingPort(key)
				: setupIncomingPort(key, callback);
		}

		managers[key] = makeManager(manager, callback);
	}

	return ports;
}

function makeManager(info, callback)
{
	var router = {
		main: callback,
		self: undefined
	};

	var tag = info.tag;
	var onEffects = info.onEffects;
	var onSelfMsg = info.onSelfMsg;

	function onMessage(msg, state)
	{
		if (msg.ctor === 'self')
		{
			return A3(onSelfMsg, router, msg._0, state);
		}

		var fx = msg._0;
		switch (tag)
		{
			case 'cmd':
				return A3(onEffects, router, fx.cmds, state);

			case 'sub':
				return A3(onEffects, router, fx.subs, state);

			case 'fx':
				return A4(onEffects, router, fx.cmds, fx.subs, state);
		}
	}

	var process = spawnLoop(info.init, onMessage);
	router.self = process;
	return process;
}

function sendToApp(router, msg)
{
	return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback)
	{
		router.main(msg);
		callback(_elm_lang$core$Native_Scheduler.succeed(_elm_lang$core$Native_Utils.Tuple0));
	});
}

function sendToSelf(router, msg)
{
	return A2(_elm_lang$core$Native_Scheduler.send, router.self, {
		ctor: 'self',
		_0: msg
	});
}


// HELPER for STATEFUL LOOPS

function spawnLoop(init, onMessage)
{
	var andThen = _elm_lang$core$Native_Scheduler.andThen;

	function loop(state)
	{
		var handleMsg = _elm_lang$core$Native_Scheduler.receive(function(msg) {
			return onMessage(msg, state);
		});
		return A2(andThen, handleMsg, loop);
	}

	var task = A2(andThen, init, loop);

	return _elm_lang$core$Native_Scheduler.rawSpawn(task);
}


// BAGS

function leaf(home)
{
	return function(value)
	{
		return {
			type: 'leaf',
			home: home,
			value: value
		};
	};
}

function batch(list)
{
	return {
		type: 'node',
		branches: list
	};
}

function map(tagger, bag)
{
	return {
		type: 'map',
		tagger: tagger,
		tree: bag
	}
}


// PIPE BAGS INTO EFFECT MANAGERS

function dispatchEffects(managers, cmdBag, subBag)
{
	var effectsDict = {};
	gatherEffects(true, cmdBag, effectsDict, null);
	gatherEffects(false, subBag, effectsDict, null);

	for (var home in managers)
	{
		var fx = home in effectsDict
			? effectsDict[home]
			: {
				cmds: _elm_lang$core$Native_List.Nil,
				subs: _elm_lang$core$Native_List.Nil
			};

		_elm_lang$core$Native_Scheduler.rawSend(managers[home], { ctor: 'fx', _0: fx });
	}
}

function gatherEffects(isCmd, bag, effectsDict, taggers)
{
	switch (bag.type)
	{
		case 'leaf':
			var home = bag.home;
			var effect = toEffect(isCmd, home, taggers, bag.value);
			effectsDict[home] = insert(isCmd, effect, effectsDict[home]);
			return;

		case 'node':
			var list = bag.branches;
			while (list.ctor !== '[]')
			{
				gatherEffects(isCmd, list._0, effectsDict, taggers);
				list = list._1;
			}
			return;

		case 'map':
			gatherEffects(isCmd, bag.tree, effectsDict, {
				tagger: bag.tagger,
				rest: taggers
			});
			return;
	}
}

function toEffect(isCmd, home, taggers, value)
{
	function applyTaggers(x)
	{
		var temp = taggers;
		while (temp)
		{
			x = temp.tagger(x);
			temp = temp.rest;
		}
		return x;
	}

	var map = isCmd
		? effectManagers[home].cmdMap
		: effectManagers[home].subMap;

	return A2(map, applyTaggers, value)
}

function insert(isCmd, newEffect, effects)
{
	effects = effects || {
		cmds: _elm_lang$core$Native_List.Nil,
		subs: _elm_lang$core$Native_List.Nil
	};
	if (isCmd)
	{
		effects.cmds = _elm_lang$core$Native_List.Cons(newEffect, effects.cmds);
		return effects;
	}
	effects.subs = _elm_lang$core$Native_List.Cons(newEffect, effects.subs);
	return effects;
}


// PORTS

function checkPortName(name)
{
	if (name in effectManagers)
	{
		throw new Error('There can only be one port named `' + name + '`, but your program has multiple.');
	}
}


// OUTGOING PORTS

function outgoingPort(name, converter)
{
	checkPortName(name);
	effectManagers[name] = {
		tag: 'cmd',
		cmdMap: outgoingPortMap,
		converter: converter,
		isForeign: true
	};
	return leaf(name);
}

var outgoingPortMap = F2(function cmdMap(tagger, value) {
	return value;
});

function setupOutgoingPort(name)
{
	var subs = [];
	var converter = effectManagers[name].converter;

	// CREATE MANAGER

	var init = _elm_lang$core$Native_Scheduler.succeed(null);

	function onEffects(router, cmdList, state)
	{
		while (cmdList.ctor !== '[]')
		{
			var value = converter(cmdList._0);
			for (var i = 0; i < subs.length; i++)
			{
				subs[i](value);
			}
			cmdList = cmdList._1;
		}
		return init;
	}

	effectManagers[name].init = init;
	effectManagers[name].onEffects = F3(onEffects);

	// PUBLIC API

	function subscribe(callback)
	{
		subs.push(callback);
	}

	function unsubscribe(callback)
	{
		var index = subs.indexOf(callback);
		if (index >= 0)
		{
			subs.splice(index, 1);
		}
	}

	return {
		subscribe: subscribe,
		unsubscribe: unsubscribe
	};
}


// INCOMING PORTS

function incomingPort(name, converter)
{
	checkPortName(name);
	effectManagers[name] = {
		tag: 'sub',
		subMap: incomingPortMap,
		converter: converter,
		isForeign: true
	};
	return leaf(name);
}

var incomingPortMap = F2(function subMap(tagger, finalTagger)
{
	return function(value)
	{
		return tagger(finalTagger(value));
	};
});

function setupIncomingPort(name, callback)
{
	var subs = _elm_lang$core$Native_List.Nil;
	var converter = effectManagers[name].converter;

	// CREATE MANAGER

	var init = _elm_lang$core$Native_Scheduler.succeed(null);

	function onEffects(router, subList, state)
	{
		subs = subList;
		return init;
	}

	effectManagers[name].init = init;
	effectManagers[name].onEffects = F3(onEffects);

	// PUBLIC API

	function send(value)
	{
		var result = A2(_elm_lang$core$Json_Decode$decodeValue, converter, value);
		if (result.ctor === 'Err')
		{
			throw new Error('Trying to send an unexpected type of value through port `' + name + '`:\n' + result._0);
		}

		var value = result._0;
		var temp = subs;
		while (temp.ctor !== '[]')
		{
			callback(temp._0(value));
			temp = temp._1;
		}
	}

	return { send: send };
}

return {
	// routers
	sendToApp: F2(sendToApp),
	sendToSelf: F2(sendToSelf),

	// global setup
	mainToProgram: mainToProgram,
	effectManagers: effectManagers,
	outgoingPort: outgoingPort,
	incomingPort: incomingPort,
	addPublicModule: addPublicModule,

	// effect bags
	leaf: leaf,
	batch: batch,
	map: F2(map)
};

}();
//import Native.Utils //

var _elm_lang$core$Native_Scheduler = function() {

var MAX_STEPS = 10000;


// TASKS

function succeed(value)
{
	return {
		ctor: '_Task_succeed',
		value: value
	};
}

function fail(error)
{
	return {
		ctor: '_Task_fail',
		value: error
	};
}

function nativeBinding(callback)
{
	return {
		ctor: '_Task_nativeBinding',
		callback: callback,
		cancel: null
	};
}

function andThen(task, callback)
{
	return {
		ctor: '_Task_andThen',
		task: task,
		callback: callback
	};
}

function onError(task, callback)
{
	return {
		ctor: '_Task_onError',
		task: task,
		callback: callback
	};
}

function receive(callback)
{
	return {
		ctor: '_Task_receive',
		callback: callback
	};
}


// PROCESSES

function rawSpawn(task)
{
	var process = {
		ctor: '_Process',
		id: _elm_lang$core$Native_Utils.guid(),
		root: task,
		stack: null,
		mailbox: []
	};

	enqueue(process);

	return process;
}

function spawn(task)
{
	return nativeBinding(function(callback) {
		var process = rawSpawn(task);
		callback(succeed(process));
	});
}

function rawSend(process, msg)
{
	process.mailbox.push(msg);
	enqueue(process);
}

function send(process, msg)
{
	return nativeBinding(function(callback) {
		rawSend(process, msg);
		callback(succeed(_elm_lang$core$Native_Utils.Tuple0));
	});
}

function kill(process)
{
	return nativeBinding(function(callback) {
		var root = process.root;
		if (root.ctor === '_Task_nativeBinding' && root.cancel)
		{
			root.cancel();
		}

		process.root = null;

		callback(succeed(_elm_lang$core$Native_Utils.Tuple0));
	});
}

function sleep(time)
{
	return nativeBinding(function(callback) {
		var id = setTimeout(function() {
			callback(succeed(_elm_lang$core$Native_Utils.Tuple0));
		}, time);

		return function() { clearTimeout(id); };
	});
}


// STEP PROCESSES

function step(numSteps, process)
{
	while (numSteps < MAX_STEPS)
	{
		var ctor = process.root.ctor;

		if (ctor === '_Task_succeed')
		{
			while (process.stack && process.stack.ctor === '_Task_onError')
			{
				process.stack = process.stack.rest;
			}
			if (process.stack === null)
			{
				break;
			}
			process.root = process.stack.callback(process.root.value);
			process.stack = process.stack.rest;
			++numSteps;
			continue;
		}

		if (ctor === '_Task_fail')
		{
			while (process.stack && process.stack.ctor === '_Task_andThen')
			{
				process.stack = process.stack.rest;
			}
			if (process.stack === null)
			{
				break;
			}
			process.root = process.stack.callback(process.root.value);
			process.stack = process.stack.rest;
			++numSteps;
			continue;
		}

		if (ctor === '_Task_andThen')
		{
			process.stack = {
				ctor: '_Task_andThen',
				callback: process.root.callback,
				rest: process.stack
			};
			process.root = process.root.task;
			++numSteps;
			continue;
		}

		if (ctor === '_Task_onError')
		{
			process.stack = {
				ctor: '_Task_onError',
				callback: process.root.callback,
				rest: process.stack
			};
			process.root = process.root.task;
			++numSteps;
			continue;
		}

		if (ctor === '_Task_nativeBinding')
		{
			process.root.cancel = process.root.callback(function(newRoot) {
				process.root = newRoot;
				enqueue(process);
			});

			break;
		}

		if (ctor === '_Task_receive')
		{
			var mailbox = process.mailbox;
			if (mailbox.length === 0)
			{
				break;
			}

			process.root = process.root.callback(mailbox.shift());
			++numSteps;
			continue;
		}

		throw new Error(ctor);
	}

	if (numSteps < MAX_STEPS)
	{
		return numSteps + 1;
	}
	enqueue(process);

	return numSteps;
}


// WORK QUEUE

var working = false;
var workQueue = [];

function enqueue(process)
{
	workQueue.push(process);

	if (!working)
	{
		setTimeout(work, 0);
		working = true;
	}
}

function work()
{
	var numSteps = 0;
	var process;
	while (numSteps < MAX_STEPS && (process = workQueue.shift()))
	{
		if (process.root)
		{
			numSteps = step(numSteps, process);
		}
	}
	if (!process)
	{
		working = false;
		return;
	}
	setTimeout(work, 0);
}


return {
	succeed: succeed,
	fail: fail,
	nativeBinding: nativeBinding,
	andThen: F2(andThen),
	onError: F2(onError),
	receive: receive,

	spawn: spawn,
	kill: kill,
	sleep: sleep,
	send: F2(send),

	rawSpawn: rawSpawn,
	rawSend: rawSend
};

}();
var _elm_lang$core$Platform$hack = _elm_lang$core$Native_Scheduler.succeed;
var _elm_lang$core$Platform$sendToSelf = _elm_lang$core$Native_Platform.sendToSelf;
var _elm_lang$core$Platform$sendToApp = _elm_lang$core$Native_Platform.sendToApp;
var _elm_lang$core$Platform$Program = {ctor: 'Program'};
var _elm_lang$core$Platform$Task = {ctor: 'Task'};
var _elm_lang$core$Platform$ProcessId = {ctor: 'ProcessId'};
var _elm_lang$core$Platform$Router = {ctor: 'Router'};

var _elm_lang$core$Platform_Cmd$batch = _elm_lang$core$Native_Platform.batch;
var _elm_lang$core$Platform_Cmd$none = _elm_lang$core$Platform_Cmd$batch(
	_elm_lang$core$Native_List.fromArray(
		[]));
var _elm_lang$core$Platform_Cmd_ops = _elm_lang$core$Platform_Cmd_ops || {};
_elm_lang$core$Platform_Cmd_ops['!'] = F2(
	function (model, commands) {
		return {
			ctor: '_Tuple2',
			_0: model,
			_1: _elm_lang$core$Platform_Cmd$batch(commands)
		};
	});
var _elm_lang$core$Platform_Cmd$map = _elm_lang$core$Native_Platform.map;
var _elm_lang$core$Platform_Cmd$Cmd = {ctor: 'Cmd'};

var _elm_lang$core$Platform_Sub$batch = _elm_lang$core$Native_Platform.batch;
var _elm_lang$core$Platform_Sub$none = _elm_lang$core$Platform_Sub$batch(
	_elm_lang$core$Native_List.fromArray(
		[]));
var _elm_lang$core$Platform_Sub$map = _elm_lang$core$Native_Platform.map;
var _elm_lang$core$Platform_Sub$Sub = {ctor: 'Sub'};

//import Native.List //

var _elm_lang$core$Native_Array = function() {

// A RRB-Tree has two distinct data types.
// Leaf -> "height"  is always 0
//         "table"   is an array of elements
// Node -> "height"  is always greater than 0
//         "table"   is an array of child nodes
//         "lengths" is an array of accumulated lengths of the child nodes

// M is the maximal table size. 32 seems fast. E is the allowed increase
// of search steps when concatting to find an index. Lower values will
// decrease balancing, but will increase search steps.
var M = 32;
var E = 2;

// An empty array.
var empty = {
	ctor: '_Array',
	height: 0,
	table: []
};


function get(i, array)
{
	if (i < 0 || i >= length(array))
	{
		throw new Error(
			'Index ' + i + ' is out of range. Check the length of ' +
			'your array first or use getMaybe or getWithDefault.');
	}
	return unsafeGet(i, array);
}


function unsafeGet(i, array)
{
	for (var x = array.height; x > 0; x--)
	{
		var slot = i >> (x * 5);
		while (array.lengths[slot] <= i)
		{
			slot++;
		}
		if (slot > 0)
		{
			i -= array.lengths[slot - 1];
		}
		array = array.table[slot];
	}
	return array.table[i];
}


// Sets the value at the index i. Only the nodes leading to i will get
// copied and updated.
function set(i, item, array)
{
	if (i < 0 || length(array) <= i)
	{
		return array;
	}
	return unsafeSet(i, item, array);
}


function unsafeSet(i, item, array)
{
	array = nodeCopy(array);

	if (array.height === 0)
	{
		array.table[i] = item;
	}
	else
	{
		var slot = getSlot(i, array);
		if (slot > 0)
		{
			i -= array.lengths[slot - 1];
		}
		array.table[slot] = unsafeSet(i, item, array.table[slot]);
	}
	return array;
}


function initialize(len, f)
{
	if (len <= 0)
	{
		return empty;
	}
	var h = Math.floor( Math.log(len) / Math.log(M) );
	return initialize_(f, h, 0, len);
}

function initialize_(f, h, from, to)
{
	if (h === 0)
	{
		var table = new Array((to - from) % (M + 1));
		for (var i = 0; i < table.length; i++)
		{
		  table[i] = f(from + i);
		}
		return {
			ctor: '_Array',
			height: 0,
			table: table
		};
	}

	var step = Math.pow(M, h);
	var table = new Array(Math.ceil((to - from) / step));
	var lengths = new Array(table.length);
	for (var i = 0; i < table.length; i++)
	{
		table[i] = initialize_(f, h - 1, from + (i * step), Math.min(from + ((i + 1) * step), to));
		lengths[i] = length(table[i]) + (i > 0 ? lengths[i-1] : 0);
	}
	return {
		ctor: '_Array',
		height: h,
		table: table,
		lengths: lengths
	};
}

function fromList(list)
{
	if (list.ctor === '[]')
	{
		return empty;
	}

	// Allocate M sized blocks (table) and write list elements to it.
	var table = new Array(M);
	var nodes = [];
	var i = 0;

	while (list.ctor !== '[]')
	{
		table[i] = list._0;
		list = list._1;
		i++;

		// table is full, so we can push a leaf containing it into the
		// next node.
		if (i === M)
		{
			var leaf = {
				ctor: '_Array',
				height: 0,
				table: table
			};
			fromListPush(leaf, nodes);
			table = new Array(M);
			i = 0;
		}
	}

	// Maybe there is something left on the table.
	if (i > 0)
	{
		var leaf = {
			ctor: '_Array',
			height: 0,
			table: table.splice(0, i)
		};
		fromListPush(leaf, nodes);
	}

	// Go through all of the nodes and eventually push them into higher nodes.
	for (var h = 0; h < nodes.length - 1; h++)
	{
		if (nodes[h].table.length > 0)
		{
			fromListPush(nodes[h], nodes);
		}
	}

	var head = nodes[nodes.length - 1];
	if (head.height > 0 && head.table.length === 1)
	{
		return head.table[0];
	}
	else
	{
		return head;
	}
}

// Push a node into a higher node as a child.
function fromListPush(toPush, nodes)
{
	var h = toPush.height;

	// Maybe the node on this height does not exist.
	if (nodes.length === h)
	{
		var node = {
			ctor: '_Array',
			height: h + 1,
			table: [],
			lengths: []
		};
		nodes.push(node);
	}

	nodes[h].table.push(toPush);
	var len = length(toPush);
	if (nodes[h].lengths.length > 0)
	{
		len += nodes[h].lengths[nodes[h].lengths.length - 1];
	}
	nodes[h].lengths.push(len);

	if (nodes[h].table.length === M)
	{
		fromListPush(nodes[h], nodes);
		nodes[h] = {
			ctor: '_Array',
			height: h + 1,
			table: [],
			lengths: []
		};
	}
}

// Pushes an item via push_ to the bottom right of a tree.
function push(item, a)
{
	var pushed = push_(item, a);
	if (pushed !== null)
	{
		return pushed;
	}

	var newTree = create(item, a.height);
	return siblise(a, newTree);
}

// Recursively tries to push an item to the bottom-right most
// tree possible. If there is no space left for the item,
// null will be returned.
function push_(item, a)
{
	// Handle resursion stop at leaf level.
	if (a.height === 0)
	{
		if (a.table.length < M)
		{
			var newA = {
				ctor: '_Array',
				height: 0,
				table: a.table.slice()
			};
			newA.table.push(item);
			return newA;
		}
		else
		{
		  return null;
		}
	}

	// Recursively push
	var pushed = push_(item, botRight(a));

	// There was space in the bottom right tree, so the slot will
	// be updated.
	if (pushed !== null)
	{
		var newA = nodeCopy(a);
		newA.table[newA.table.length - 1] = pushed;
		newA.lengths[newA.lengths.length - 1]++;
		return newA;
	}

	// When there was no space left, check if there is space left
	// for a new slot with a tree which contains only the item
	// at the bottom.
	if (a.table.length < M)
	{
		var newSlot = create(item, a.height - 1);
		var newA = nodeCopy(a);
		newA.table.push(newSlot);
		newA.lengths.push(newA.lengths[newA.lengths.length - 1] + length(newSlot));
		return newA;
	}
	else
	{
		return null;
	}
}

// Converts an array into a list of elements.
function toList(a)
{
	return toList_(_elm_lang$core$Native_List.Nil, a);
}

function toList_(list, a)
{
	for (var i = a.table.length - 1; i >= 0; i--)
	{
		list =
			a.height === 0
				? _elm_lang$core$Native_List.Cons(a.table[i], list)
				: toList_(list, a.table[i]);
	}
	return list;
}

// Maps a function over the elements of an array.
function map(f, a)
{
	var newA = {
		ctor: '_Array',
		height: a.height,
		table: new Array(a.table.length)
	};
	if (a.height > 0)
	{
		newA.lengths = a.lengths;
	}
	for (var i = 0; i < a.table.length; i++)
	{
		newA.table[i] =
			a.height === 0
				? f(a.table[i])
				: map(f, a.table[i]);
	}
	return newA;
}

// Maps a function over the elements with their index as first argument.
function indexedMap(f, a)
{
	return indexedMap_(f, a, 0);
}

function indexedMap_(f, a, from)
{
	var newA = {
		ctor: '_Array',
		height: a.height,
		table: new Array(a.table.length)
	};
	if (a.height > 0)
	{
		newA.lengths = a.lengths;
	}
	for (var i = 0; i < a.table.length; i++)
	{
		newA.table[i] =
			a.height === 0
				? A2(f, from + i, a.table[i])
				: indexedMap_(f, a.table[i], i == 0 ? from : from + a.lengths[i - 1]);
	}
	return newA;
}

function foldl(f, b, a)
{
	if (a.height === 0)
	{
		for (var i = 0; i < a.table.length; i++)
		{
			b = A2(f, a.table[i], b);
		}
	}
	else
	{
		for (var i = 0; i < a.table.length; i++)
		{
			b = foldl(f, b, a.table[i]);
		}
	}
	return b;
}

function foldr(f, b, a)
{
	if (a.height === 0)
	{
		for (var i = a.table.length; i--; )
		{
			b = A2(f, a.table[i], b);
		}
	}
	else
	{
		for (var i = a.table.length; i--; )
		{
			b = foldr(f, b, a.table[i]);
		}
	}
	return b;
}

// TODO: currently, it slices the right, then the left. This can be
// optimized.
function slice(from, to, a)
{
	if (from < 0)
	{
		from += length(a);
	}
	if (to < 0)
	{
		to += length(a);
	}
	return sliceLeft(from, sliceRight(to, a));
}

function sliceRight(to, a)
{
	if (to === length(a))
	{
		return a;
	}

	// Handle leaf level.
	if (a.height === 0)
	{
		var newA = { ctor:'_Array', height:0 };
		newA.table = a.table.slice(0, to);
		return newA;
	}

	// Slice the right recursively.
	var right = getSlot(to, a);
	var sliced = sliceRight(to - (right > 0 ? a.lengths[right - 1] : 0), a.table[right]);

	// Maybe the a node is not even needed, as sliced contains the whole slice.
	if (right === 0)
	{
		return sliced;
	}

	// Create new node.
	var newA = {
		ctor: '_Array',
		height: a.height,
		table: a.table.slice(0, right),
		lengths: a.lengths.slice(0, right)
	};
	if (sliced.table.length > 0)
	{
		newA.table[right] = sliced;
		newA.lengths[right] = length(sliced) + (right > 0 ? newA.lengths[right - 1] : 0);
	}
	return newA;
}

function sliceLeft(from, a)
{
	if (from === 0)
	{
		return a;
	}

	// Handle leaf level.
	if (a.height === 0)
	{
		var newA = { ctor:'_Array', height:0 };
		newA.table = a.table.slice(from, a.table.length + 1);
		return newA;
	}

	// Slice the left recursively.
	var left = getSlot(from, a);
	var sliced = sliceLeft(from - (left > 0 ? a.lengths[left - 1] : 0), a.table[left]);

	// Maybe the a node is not even needed, as sliced contains the whole slice.
	if (left === a.table.length - 1)
	{
		return sliced;
	}

	// Create new node.
	var newA = {
		ctor: '_Array',
		height: a.height,
		table: a.table.slice(left, a.table.length + 1),
		lengths: new Array(a.table.length - left)
	};
	newA.table[0] = sliced;
	var len = 0;
	for (var i = 0; i < newA.table.length; i++)
	{
		len += length(newA.table[i]);
		newA.lengths[i] = len;
	}

	return newA;
}

// Appends two trees.
function append(a,b)
{
	if (a.table.length === 0)
	{
		return b;
	}
	if (b.table.length === 0)
	{
		return a;
	}

	var c = append_(a, b);

	// Check if both nodes can be crunshed together.
	if (c[0].table.length + c[1].table.length <= M)
	{
		if (c[0].table.length === 0)
		{
			return c[1];
		}
		if (c[1].table.length === 0)
		{
			return c[0];
		}

		// Adjust .table and .lengths
		c[0].table = c[0].table.concat(c[1].table);
		if (c[0].height > 0)
		{
			var len = length(c[0]);
			for (var i = 0; i < c[1].lengths.length; i++)
			{
				c[1].lengths[i] += len;
			}
			c[0].lengths = c[0].lengths.concat(c[1].lengths);
		}

		return c[0];
	}

	if (c[0].height > 0)
	{
		var toRemove = calcToRemove(a, b);
		if (toRemove > E)
		{
			c = shuffle(c[0], c[1], toRemove);
		}
	}

	return siblise(c[0], c[1]);
}

// Returns an array of two nodes; right and left. One node _may_ be empty.
function append_(a, b)
{
	if (a.height === 0 && b.height === 0)
	{
		return [a, b];
	}

	if (a.height !== 1 || b.height !== 1)
	{
		if (a.height === b.height)
		{
			a = nodeCopy(a);
			b = nodeCopy(b);
			var appended = append_(botRight(a), botLeft(b));

			insertRight(a, appended[1]);
			insertLeft(b, appended[0]);
		}
		else if (a.height > b.height)
		{
			a = nodeCopy(a);
			var appended = append_(botRight(a), b);

			insertRight(a, appended[0]);
			b = parentise(appended[1], appended[1].height + 1);
		}
		else
		{
			b = nodeCopy(b);
			var appended = append_(a, botLeft(b));

			var left = appended[0].table.length === 0 ? 0 : 1;
			var right = left === 0 ? 1 : 0;
			insertLeft(b, appended[left]);
			a = parentise(appended[right], appended[right].height + 1);
		}
	}

	// Check if balancing is needed and return based on that.
	if (a.table.length === 0 || b.table.length === 0)
	{
		return [a, b];
	}

	var toRemove = calcToRemove(a, b);
	if (toRemove <= E)
	{
		return [a, b];
	}
	return shuffle(a, b, toRemove);
}

// Helperfunctions for append_. Replaces a child node at the side of the parent.
function insertRight(parent, node)
{
	var index = parent.table.length - 1;
	parent.table[index] = node;
	parent.lengths[index] = length(node);
	parent.lengths[index] += index > 0 ? parent.lengths[index - 1] : 0;
}

function insertLeft(parent, node)
{
	if (node.table.length > 0)
	{
		parent.table[0] = node;
		parent.lengths[0] = length(node);

		var len = length(parent.table[0]);
		for (var i = 1; i < parent.lengths.length; i++)
		{
			len += length(parent.table[i]);
			parent.lengths[i] = len;
		}
	}
	else
	{
		parent.table.shift();
		for (var i = 1; i < parent.lengths.length; i++)
		{
			parent.lengths[i] = parent.lengths[i] - parent.lengths[0];
		}
		parent.lengths.shift();
	}
}

// Returns the extra search steps for E. Refer to the paper.
function calcToRemove(a, b)
{
	var subLengths = 0;
	for (var i = 0; i < a.table.length; i++)
	{
		subLengths += a.table[i].table.length;
	}
	for (var i = 0; i < b.table.length; i++)
	{
		subLengths += b.table[i].table.length;
	}

	var toRemove = a.table.length + b.table.length;
	return toRemove - (Math.floor((subLengths - 1) / M) + 1);
}

// get2, set2 and saveSlot are helpers for accessing elements over two arrays.
function get2(a, b, index)
{
	return index < a.length
		? a[index]
		: b[index - a.length];
}

function set2(a, b, index, value)
{
	if (index < a.length)
	{
		a[index] = value;
	}
	else
	{
		b[index - a.length] = value;
	}
}

function saveSlot(a, b, index, slot)
{
	set2(a.table, b.table, index, slot);

	var l = (index === 0 || index === a.lengths.length)
		? 0
		: get2(a.lengths, a.lengths, index - 1);

	set2(a.lengths, b.lengths, index, l + length(slot));
}

// Creates a node or leaf with a given length at their arrays for perfomance.
// Is only used by shuffle.
function createNode(h, length)
{
	if (length < 0)
	{
		length = 0;
	}
	var a = {
		ctor: '_Array',
		height: h,
		table: new Array(length)
	};
	if (h > 0)
	{
		a.lengths = new Array(length);
	}
	return a;
}

// Returns an array of two balanced nodes.
function shuffle(a, b, toRemove)
{
	var newA = createNode(a.height, Math.min(M, a.table.length + b.table.length - toRemove));
	var newB = createNode(a.height, newA.table.length - (a.table.length + b.table.length - toRemove));

	// Skip the slots with size M. More precise: copy the slot references
	// to the new node
	var read = 0;
	while (get2(a.table, b.table, read).table.length % M === 0)
	{
		set2(newA.table, newB.table, read, get2(a.table, b.table, read));
		set2(newA.lengths, newB.lengths, read, get2(a.lengths, b.lengths, read));
		read++;
	}

	// Pulling items from left to right, caching in a slot before writing
	// it into the new nodes.
	var write = read;
	var slot = new createNode(a.height - 1, 0);
	var from = 0;

	// If the current slot is still containing data, then there will be at
	// least one more write, so we do not break this loop yet.
	while (read - write - (slot.table.length > 0 ? 1 : 0) < toRemove)
	{
		// Find out the max possible items for copying.
		var source = get2(a.table, b.table, read);
		var to = Math.min(M - slot.table.length, source.table.length);

		// Copy and adjust size table.
		slot.table = slot.table.concat(source.table.slice(from, to));
		if (slot.height > 0)
		{
			var len = slot.lengths.length;
			for (var i = len; i < len + to - from; i++)
			{
				slot.lengths[i] = length(slot.table[i]);
				slot.lengths[i] += (i > 0 ? slot.lengths[i - 1] : 0);
			}
		}

		from += to;

		// Only proceed to next slots[i] if the current one was
		// fully copied.
		if (source.table.length <= to)
		{
			read++; from = 0;
		}

		// Only create a new slot if the current one is filled up.
		if (slot.table.length === M)
		{
			saveSlot(newA, newB, write, slot);
			slot = createNode(a.height - 1, 0);
			write++;
		}
	}

	// Cleanup after the loop. Copy the last slot into the new nodes.
	if (slot.table.length > 0)
	{
		saveSlot(newA, newB, write, slot);
		write++;
	}

	// Shift the untouched slots to the left
	while (read < a.table.length + b.table.length )
	{
		saveSlot(newA, newB, write, get2(a.table, b.table, read));
		read++;
		write++;
	}

	return [newA, newB];
}

// Navigation functions
function botRight(a)
{
	return a.table[a.table.length - 1];
}
function botLeft(a)
{
	return a.table[0];
}

// Copies a node for updating. Note that you should not use this if
// only updating only one of "table" or "lengths" for performance reasons.
function nodeCopy(a)
{
	var newA = {
		ctor: '_Array',
		height: a.height,
		table: a.table.slice()
	};
	if (a.height > 0)
	{
		newA.lengths = a.lengths.slice();
	}
	return newA;
}

// Returns how many items are in the tree.
function length(array)
{
	if (array.height === 0)
	{
		return array.table.length;
	}
	else
	{
		return array.lengths[array.lengths.length - 1];
	}
}

// Calculates in which slot of "table" the item probably is, then
// find the exact slot via forward searching in  "lengths". Returns the index.
function getSlot(i, a)
{
	var slot = i >> (5 * a.height);
	while (a.lengths[slot] <= i)
	{
		slot++;
	}
	return slot;
}

// Recursively creates a tree with a given height containing
// only the given item.
function create(item, h)
{
	if (h === 0)
	{
		return {
			ctor: '_Array',
			height: 0,
			table: [item]
		};
	}
	return {
		ctor: '_Array',
		height: h,
		table: [create(item, h - 1)],
		lengths: [1]
	};
}

// Recursively creates a tree that contains the given tree.
function parentise(tree, h)
{
	if (h === tree.height)
	{
		return tree;
	}

	return {
		ctor: '_Array',
		height: h,
		table: [parentise(tree, h - 1)],
		lengths: [length(tree)]
	};
}

// Emphasizes blood brotherhood beneath two trees.
function siblise(a, b)
{
	return {
		ctor: '_Array',
		height: a.height + 1,
		table: [a, b],
		lengths: [length(a), length(a) + length(b)]
	};
}

function toJSArray(a)
{
	var jsArray = new Array(length(a));
	toJSArray_(jsArray, 0, a);
	return jsArray;
}

function toJSArray_(jsArray, i, a)
{
	for (var t = 0; t < a.table.length; t++)
	{
		if (a.height === 0)
		{
			jsArray[i + t] = a.table[t];
		}
		else
		{
			var inc = t === 0 ? 0 : a.lengths[t - 1];
			toJSArray_(jsArray, i + inc, a.table[t]);
		}
	}
}

function fromJSArray(jsArray)
{
	if (jsArray.length === 0)
	{
		return empty;
	}
	var h = Math.floor(Math.log(jsArray.length) / Math.log(M));
	return fromJSArray_(jsArray, h, 0, jsArray.length);
}

function fromJSArray_(jsArray, h, from, to)
{
	if (h === 0)
	{
		return {
			ctor: '_Array',
			height: 0,
			table: jsArray.slice(from, to)
		};
	}

	var step = Math.pow(M, h);
	var table = new Array(Math.ceil((to - from) / step));
	var lengths = new Array(table.length);
	for (var i = 0; i < table.length; i++)
	{
		table[i] = fromJSArray_(jsArray, h - 1, from + (i * step), Math.min(from + ((i + 1) * step), to));
		lengths[i] = length(table[i]) + (i > 0 ? lengths[i - 1] : 0);
	}
	return {
		ctor: '_Array',
		height: h,
		table: table,
		lengths: lengths
	};
}

return {
	empty: empty,
	fromList: fromList,
	toList: toList,
	initialize: F2(initialize),
	append: F2(append),
	push: F2(push),
	slice: F3(slice),
	get: F2(get),
	set: F3(set),
	map: F2(map),
	indexedMap: F2(indexedMap),
	foldl: F3(foldl),
	foldr: F3(foldr),
	length: length,

	toJSArray: toJSArray,
	fromJSArray: fromJSArray
};

}();
var _elm_lang$core$Array$append = _elm_lang$core$Native_Array.append;
var _elm_lang$core$Array$length = _elm_lang$core$Native_Array.length;
var _elm_lang$core$Array$isEmpty = function (array) {
	return _elm_lang$core$Native_Utils.eq(
		_elm_lang$core$Array$length(array),
		0);
};
var _elm_lang$core$Array$slice = _elm_lang$core$Native_Array.slice;
var _elm_lang$core$Array$set = _elm_lang$core$Native_Array.set;
var _elm_lang$core$Array$get = F2(
	function (i, array) {
		return ((_elm_lang$core$Native_Utils.cmp(0, i) < 1) && (_elm_lang$core$Native_Utils.cmp(
			i,
			_elm_lang$core$Native_Array.length(array)) < 0)) ? _elm_lang$core$Maybe$Just(
			A2(_elm_lang$core$Native_Array.get, i, array)) : _elm_lang$core$Maybe$Nothing;
	});
var _elm_lang$core$Array$push = _elm_lang$core$Native_Array.push;
var _elm_lang$core$Array$empty = _elm_lang$core$Native_Array.empty;
var _elm_lang$core$Array$filter = F2(
	function (isOkay, arr) {
		var update = F2(
			function (x, xs) {
				return isOkay(x) ? A2(_elm_lang$core$Native_Array.push, x, xs) : xs;
			});
		return A3(_elm_lang$core$Native_Array.foldl, update, _elm_lang$core$Native_Array.empty, arr);
	});
var _elm_lang$core$Array$foldr = _elm_lang$core$Native_Array.foldr;
var _elm_lang$core$Array$foldl = _elm_lang$core$Native_Array.foldl;
var _elm_lang$core$Array$indexedMap = _elm_lang$core$Native_Array.indexedMap;
var _elm_lang$core$Array$map = _elm_lang$core$Native_Array.map;
var _elm_lang$core$Array$toIndexedList = function (array) {
	return A3(
		_elm_lang$core$List$map2,
		F2(
			function (v0, v1) {
				return {ctor: '_Tuple2', _0: v0, _1: v1};
			}),
		_elm_lang$core$Native_List.range(
			0,
			_elm_lang$core$Native_Array.length(array) - 1),
		_elm_lang$core$Native_Array.toList(array));
};
var _elm_lang$core$Array$toList = _elm_lang$core$Native_Array.toList;
var _elm_lang$core$Array$fromList = _elm_lang$core$Native_Array.fromList;
var _elm_lang$core$Array$initialize = _elm_lang$core$Native_Array.initialize;
var _elm_lang$core$Array$repeat = F2(
	function (n, e) {
		return A2(
			_elm_lang$core$Array$initialize,
			n,
			_elm_lang$core$Basics$always(e));
	});
var _elm_lang$core$Array$Array = {ctor: 'Array'};

var _chendrix$elm_matrix$Matrix$rowCount = function (m) {
	return _elm_lang$core$Array$length(m);
};
var _chendrix$elm_matrix$Matrix$colCount = function (m) {
	return A2(
		_elm_lang$core$Maybe$withDefault,
		0,
		A2(
			_elm_lang$core$Maybe$map,
			_elm_lang$core$Array$length,
			A2(_elm_lang$core$Array$get, 0, m)));
};
var _chendrix$elm_matrix$Matrix$fromList = function (l) {
	return _elm_lang$core$Array$fromList(
		A2(_elm_lang$core$List$map, _elm_lang$core$Array$fromList, l));
};
var _chendrix$elm_matrix$Matrix$toList = function (m) {
	return _elm_lang$core$Array$toList(
		A2(_elm_lang$core$Array$map, _elm_lang$core$Array$toList, m));
};
var _chendrix$elm_matrix$Matrix$flatten = function (m) {
	return _elm_lang$core$List$concat(
		_chendrix$elm_matrix$Matrix$toList(m));
};
var _chendrix$elm_matrix$Matrix$map = F2(
	function (f, m) {
		return A2(
			_elm_lang$core$Array$map,
			_elm_lang$core$Array$map(f),
			m);
	});
var _chendrix$elm_matrix$Matrix$col = _elm_lang$core$Basics$snd;
var _chendrix$elm_matrix$Matrix$row = _elm_lang$core$Basics$fst;
var _chendrix$elm_matrix$Matrix$get = F2(
	function (location, m) {
		return A2(
			_elm_lang$core$Maybe$andThen,
			A2(
				_elm_lang$core$Array$get,
				_chendrix$elm_matrix$Matrix$row(location),
				m),
			_elm_lang$core$Array$get(
				_chendrix$elm_matrix$Matrix$col(location)));
	});
var _chendrix$elm_matrix$Matrix$update = F3(
	function (location, f, m) {
		return A2(
			_elm_lang$core$Maybe$withDefault,
			m,
			A2(
				_elm_lang$core$Maybe$map,
				function (current) {
					return A2(
						_elm_lang$core$Maybe$withDefault,
						m,
						A2(
							_elm_lang$core$Maybe$map,
							function (oldRow) {
								return function (newRow) {
									return A3(
										_elm_lang$core$Array$set,
										_chendrix$elm_matrix$Matrix$row(location),
										newRow,
										m);
								}(
									A3(
										_elm_lang$core$Array$set,
										_chendrix$elm_matrix$Matrix$col(location),
										f(current),
										oldRow));
							},
							A2(
								_elm_lang$core$Array$get,
								_chendrix$elm_matrix$Matrix$row(location),
								m)));
				},
				A2(_chendrix$elm_matrix$Matrix$get, location, m)));
	});
var _chendrix$elm_matrix$Matrix$set = F3(
	function (location, value, m) {
		return A3(
			_chendrix$elm_matrix$Matrix$update,
			location,
			_elm_lang$core$Basics$always(value),
			m);
	});
var _chendrix$elm_matrix$Matrix$loc = F2(
	function (v0, v1) {
		return {ctor: '_Tuple2', _0: v0, _1: v1};
	});
var _chendrix$elm_matrix$Matrix$matrix = F3(
	function (numRows, numCols, f) {
		return A2(
			_elm_lang$core$Array$initialize,
			numRows,
			function (row) {
				return A2(
					_elm_lang$core$Array$initialize,
					numCols,
					function (col) {
						return f(
							A2(_chendrix$elm_matrix$Matrix$loc, row, col));
					});
			});
	});
var _chendrix$elm_matrix$Matrix$square = function (size) {
	return A2(_chendrix$elm_matrix$Matrix$matrix, size, size);
};
var _chendrix$elm_matrix$Matrix$mapWithLocation = F2(
	function (f, m) {
		return A2(
			_elm_lang$core$Array$indexedMap,
			F2(
				function (rowNum, row) {
					return A2(
						_elm_lang$core$Array$indexedMap,
						F2(
							function (colNum, element) {
								return A2(
									f,
									A2(_chendrix$elm_matrix$Matrix$loc, rowNum, colNum),
									element);
							}),
						row);
				}),
			m);
	});

var _elm_lang$core$Task$onError = _elm_lang$core$Native_Scheduler.onError;
var _elm_lang$core$Task$andThen = _elm_lang$core$Native_Scheduler.andThen;
var _elm_lang$core$Task$spawnCmd = F2(
	function (router, _p0) {
		var _p1 = _p0;
		return _elm_lang$core$Native_Scheduler.spawn(
			A2(
				_elm_lang$core$Task$andThen,
				_p1._0,
				_elm_lang$core$Platform$sendToApp(router)));
	});
var _elm_lang$core$Task$fail = _elm_lang$core$Native_Scheduler.fail;
var _elm_lang$core$Task$mapError = F2(
	function (f, task) {
		return A2(
			_elm_lang$core$Task$onError,
			task,
			function (err) {
				return _elm_lang$core$Task$fail(
					f(err));
			});
	});
var _elm_lang$core$Task$succeed = _elm_lang$core$Native_Scheduler.succeed;
var _elm_lang$core$Task$map = F2(
	function (func, taskA) {
		return A2(
			_elm_lang$core$Task$andThen,
			taskA,
			function (a) {
				return _elm_lang$core$Task$succeed(
					func(a));
			});
	});
var _elm_lang$core$Task$map2 = F3(
	function (func, taskA, taskB) {
		return A2(
			_elm_lang$core$Task$andThen,
			taskA,
			function (a) {
				return A2(
					_elm_lang$core$Task$andThen,
					taskB,
					function (b) {
						return _elm_lang$core$Task$succeed(
							A2(func, a, b));
					});
			});
	});
var _elm_lang$core$Task$map3 = F4(
	function (func, taskA, taskB, taskC) {
		return A2(
			_elm_lang$core$Task$andThen,
			taskA,
			function (a) {
				return A2(
					_elm_lang$core$Task$andThen,
					taskB,
					function (b) {
						return A2(
							_elm_lang$core$Task$andThen,
							taskC,
							function (c) {
								return _elm_lang$core$Task$succeed(
									A3(func, a, b, c));
							});
					});
			});
	});
var _elm_lang$core$Task$map4 = F5(
	function (func, taskA, taskB, taskC, taskD) {
		return A2(
			_elm_lang$core$Task$andThen,
			taskA,
			function (a) {
				return A2(
					_elm_lang$core$Task$andThen,
					taskB,
					function (b) {
						return A2(
							_elm_lang$core$Task$andThen,
							taskC,
							function (c) {
								return A2(
									_elm_lang$core$Task$andThen,
									taskD,
									function (d) {
										return _elm_lang$core$Task$succeed(
											A4(func, a, b, c, d));
									});
							});
					});
			});
	});
var _elm_lang$core$Task$map5 = F6(
	function (func, taskA, taskB, taskC, taskD, taskE) {
		return A2(
			_elm_lang$core$Task$andThen,
			taskA,
			function (a) {
				return A2(
					_elm_lang$core$Task$andThen,
					taskB,
					function (b) {
						return A2(
							_elm_lang$core$Task$andThen,
							taskC,
							function (c) {
								return A2(
									_elm_lang$core$Task$andThen,
									taskD,
									function (d) {
										return A2(
											_elm_lang$core$Task$andThen,
											taskE,
											function (e) {
												return _elm_lang$core$Task$succeed(
													A5(func, a, b, c, d, e));
											});
									});
							});
					});
			});
	});
var _elm_lang$core$Task$andMap = F2(
	function (taskFunc, taskValue) {
		return A2(
			_elm_lang$core$Task$andThen,
			taskFunc,
			function (func) {
				return A2(
					_elm_lang$core$Task$andThen,
					taskValue,
					function (value) {
						return _elm_lang$core$Task$succeed(
							func(value));
					});
			});
	});
var _elm_lang$core$Task$sequence = function (tasks) {
	var _p2 = tasks;
	if (_p2.ctor === '[]') {
		return _elm_lang$core$Task$succeed(
			_elm_lang$core$Native_List.fromArray(
				[]));
	} else {
		return A3(
			_elm_lang$core$Task$map2,
			F2(
				function (x, y) {
					return A2(_elm_lang$core$List_ops['::'], x, y);
				}),
			_p2._0,
			_elm_lang$core$Task$sequence(_p2._1));
	}
};
var _elm_lang$core$Task$onEffects = F3(
	function (router, commands, state) {
		return A2(
			_elm_lang$core$Task$map,
			function (_p3) {
				return {ctor: '_Tuple0'};
			},
			_elm_lang$core$Task$sequence(
				A2(
					_elm_lang$core$List$map,
					_elm_lang$core$Task$spawnCmd(router),
					commands)));
	});
var _elm_lang$core$Task$toMaybe = function (task) {
	return A2(
		_elm_lang$core$Task$onError,
		A2(_elm_lang$core$Task$map, _elm_lang$core$Maybe$Just, task),
		function (_p4) {
			return _elm_lang$core$Task$succeed(_elm_lang$core$Maybe$Nothing);
		});
};
var _elm_lang$core$Task$fromMaybe = F2(
	function ($default, maybe) {
		var _p5 = maybe;
		if (_p5.ctor === 'Just') {
			return _elm_lang$core$Task$succeed(_p5._0);
		} else {
			return _elm_lang$core$Task$fail($default);
		}
	});
var _elm_lang$core$Task$toResult = function (task) {
	return A2(
		_elm_lang$core$Task$onError,
		A2(_elm_lang$core$Task$map, _elm_lang$core$Result$Ok, task),
		function (msg) {
			return _elm_lang$core$Task$succeed(
				_elm_lang$core$Result$Err(msg));
		});
};
var _elm_lang$core$Task$fromResult = function (result) {
	var _p6 = result;
	if (_p6.ctor === 'Ok') {
		return _elm_lang$core$Task$succeed(_p6._0);
	} else {
		return _elm_lang$core$Task$fail(_p6._0);
	}
};
var _elm_lang$core$Task$init = _elm_lang$core$Task$succeed(
	{ctor: '_Tuple0'});
var _elm_lang$core$Task$onSelfMsg = F3(
	function (_p9, _p8, _p7) {
		return _elm_lang$core$Task$succeed(
			{ctor: '_Tuple0'});
	});
var _elm_lang$core$Task$command = _elm_lang$core$Native_Platform.leaf('Task');
var _elm_lang$core$Task$T = function (a) {
	return {ctor: 'T', _0: a};
};
var _elm_lang$core$Task$perform = F3(
	function (onFail, onSuccess, task) {
		return _elm_lang$core$Task$command(
			_elm_lang$core$Task$T(
				A2(
					_elm_lang$core$Task$onError,
					A2(_elm_lang$core$Task$map, onSuccess, task),
					function (x) {
						return _elm_lang$core$Task$succeed(
							onFail(x));
					})));
	});
var _elm_lang$core$Task$cmdMap = F2(
	function (tagger, _p10) {
		var _p11 = _p10;
		return _elm_lang$core$Task$T(
			A2(_elm_lang$core$Task$map, tagger, _p11._0));
	});
_elm_lang$core$Native_Platform.effectManagers['Task'] = {pkg: 'elm-lang/core', init: _elm_lang$core$Task$init, onEffects: _elm_lang$core$Task$onEffects, onSelfMsg: _elm_lang$core$Task$onSelfMsg, tag: 'cmd', cmdMap: _elm_lang$core$Task$cmdMap};

//import Maybe, Native.List, Native.Utils, Result //

var _elm_lang$core$Native_String = function() {

function isEmpty(str)
{
	return str.length === 0;
}
function cons(chr, str)
{
	return chr + str;
}
function uncons(str)
{
	var hd = str[0];
	if (hd)
	{
		return _elm_lang$core$Maybe$Just(_elm_lang$core$Native_Utils.Tuple2(_elm_lang$core$Native_Utils.chr(hd), str.slice(1)));
	}
	return _elm_lang$core$Maybe$Nothing;
}
function append(a, b)
{
	return a + b;
}
function concat(strs)
{
	return _elm_lang$core$Native_List.toArray(strs).join('');
}
function length(str)
{
	return str.length;
}
function map(f, str)
{
	var out = str.split('');
	for (var i = out.length; i--; )
	{
		out[i] = f(_elm_lang$core$Native_Utils.chr(out[i]));
	}
	return out.join('');
}
function filter(pred, str)
{
	return str.split('').map(_elm_lang$core$Native_Utils.chr).filter(pred).join('');
}
function reverse(str)
{
	return str.split('').reverse().join('');
}
function foldl(f, b, str)
{
	var len = str.length;
	for (var i = 0; i < len; ++i)
	{
		b = A2(f, _elm_lang$core$Native_Utils.chr(str[i]), b);
	}
	return b;
}
function foldr(f, b, str)
{
	for (var i = str.length; i--; )
	{
		b = A2(f, _elm_lang$core$Native_Utils.chr(str[i]), b);
	}
	return b;
}
function split(sep, str)
{
	return _elm_lang$core$Native_List.fromArray(str.split(sep));
}
function join(sep, strs)
{
	return _elm_lang$core$Native_List.toArray(strs).join(sep);
}
function repeat(n, str)
{
	var result = '';
	while (n > 0)
	{
		if (n & 1)
		{
			result += str;
		}
		n >>= 1, str += str;
	}
	return result;
}
function slice(start, end, str)
{
	return str.slice(start, end);
}
function left(n, str)
{
	return n < 1 ? '' : str.slice(0, n);
}
function right(n, str)
{
	return n < 1 ? '' : str.slice(-n);
}
function dropLeft(n, str)
{
	return n < 1 ? str : str.slice(n);
}
function dropRight(n, str)
{
	return n < 1 ? str : str.slice(0, -n);
}
function pad(n, chr, str)
{
	var half = (n - str.length) / 2;
	return repeat(Math.ceil(half), chr) + str + repeat(half | 0, chr);
}
function padRight(n, chr, str)
{
	return str + repeat(n - str.length, chr);
}
function padLeft(n, chr, str)
{
	return repeat(n - str.length, chr) + str;
}

function trim(str)
{
	return str.trim();
}
function trimLeft(str)
{
	return str.replace(/^\s+/, '');
}
function trimRight(str)
{
	return str.replace(/\s+$/, '');
}

function words(str)
{
	return _elm_lang$core$Native_List.fromArray(str.trim().split(/\s+/g));
}
function lines(str)
{
	return _elm_lang$core$Native_List.fromArray(str.split(/\r\n|\r|\n/g));
}

function toUpper(str)
{
	return str.toUpperCase();
}
function toLower(str)
{
	return str.toLowerCase();
}

function any(pred, str)
{
	for (var i = str.length; i--; )
	{
		if (pred(_elm_lang$core$Native_Utils.chr(str[i])))
		{
			return true;
		}
	}
	return false;
}
function all(pred, str)
{
	for (var i = str.length; i--; )
	{
		if (!pred(_elm_lang$core$Native_Utils.chr(str[i])))
		{
			return false;
		}
	}
	return true;
}

function contains(sub, str)
{
	return str.indexOf(sub) > -1;
}
function startsWith(sub, str)
{
	return str.indexOf(sub) === 0;
}
function endsWith(sub, str)
{
	return str.length >= sub.length &&
		str.lastIndexOf(sub) === str.length - sub.length;
}
function indexes(sub, str)
{
	var subLen = sub.length;
	
	if (subLen < 1)
	{
		return _elm_lang$core$Native_List.Nil;
	}

	var i = 0;
	var is = [];

	while ((i = str.indexOf(sub, i)) > -1)
	{
		is.push(i);
		i = i + subLen;
	}	
	
	return _elm_lang$core$Native_List.fromArray(is);
}

function toInt(s)
{
	var len = s.length;
	if (len === 0)
	{
		return _elm_lang$core$Result$Err("could not convert string '" + s + "' to an Int" );
	}
	var start = 0;
	if (s[0] === '-')
	{
		if (len === 1)
		{
			return _elm_lang$core$Result$Err("could not convert string '" + s + "' to an Int" );
		}
		start = 1;
	}
	for (var i = start; i < len; ++i)
	{
		var c = s[i];
		if (c < '0' || '9' < c)
		{
			return _elm_lang$core$Result$Err("could not convert string '" + s + "' to an Int" );
		}
	}
	return _elm_lang$core$Result$Ok(parseInt(s, 10));
}

function toFloat(s)
{
	var len = s.length;
	if (len === 0)
	{
		return _elm_lang$core$Result$Err("could not convert string '" + s + "' to a Float" );
	}
	var start = 0;
	if (s[0] === '-')
	{
		if (len === 1)
		{
			return _elm_lang$core$Result$Err("could not convert string '" + s + "' to a Float" );
		}
		start = 1;
	}
	var dotCount = 0;
	for (var i = start; i < len; ++i)
	{
		var c = s[i];
		if ('0' <= c && c <= '9')
		{
			continue;
		}
		if (c === '.')
		{
			dotCount += 1;
			if (dotCount <= 1)
			{
				continue;
			}
		}
		return _elm_lang$core$Result$Err("could not convert string '" + s + "' to a Float" );
	}
	return _elm_lang$core$Result$Ok(parseFloat(s));
}

function toList(str)
{
	return _elm_lang$core$Native_List.fromArray(str.split('').map(_elm_lang$core$Native_Utils.chr));
}
function fromList(chars)
{
	return _elm_lang$core$Native_List.toArray(chars).join('');
}

return {
	isEmpty: isEmpty,
	cons: F2(cons),
	uncons: uncons,
	append: F2(append),
	concat: concat,
	length: length,
	map: F2(map),
	filter: F2(filter),
	reverse: reverse,
	foldl: F3(foldl),
	foldr: F3(foldr),

	split: F2(split),
	join: F2(join),
	repeat: F2(repeat),

	slice: F3(slice),
	left: F2(left),
	right: F2(right),
	dropLeft: F2(dropLeft),
	dropRight: F2(dropRight),

	pad: F3(pad),
	padLeft: F3(padLeft),
	padRight: F3(padRight),

	trim: trim,
	trimLeft: trimLeft,
	trimRight: trimRight,

	words: words,
	lines: lines,

	toUpper: toUpper,
	toLower: toLower,

	any: F2(any),
	all: F2(all),

	contains: F2(contains),
	startsWith: F2(startsWith),
	endsWith: F2(endsWith),
	indexes: F2(indexes),

	toInt: toInt,
	toFloat: toFloat,
	toList: toList,
	fromList: fromList
};

}();

//import Native.Utils //

var _elm_lang$core$Native_Char = function() {

return {
	fromCode: function(c) { return _elm_lang$core$Native_Utils.chr(String.fromCharCode(c)); },
	toCode: function(c) { return c.charCodeAt(0); },
	toUpper: function(c) { return _elm_lang$core$Native_Utils.chr(c.toUpperCase()); },
	toLower: function(c) { return _elm_lang$core$Native_Utils.chr(c.toLowerCase()); },
	toLocaleUpper: function(c) { return _elm_lang$core$Native_Utils.chr(c.toLocaleUpperCase()); },
	toLocaleLower: function(c) { return _elm_lang$core$Native_Utils.chr(c.toLocaleLowerCase()); }
};

}();
var _elm_lang$core$Char$fromCode = _elm_lang$core$Native_Char.fromCode;
var _elm_lang$core$Char$toCode = _elm_lang$core$Native_Char.toCode;
var _elm_lang$core$Char$toLocaleLower = _elm_lang$core$Native_Char.toLocaleLower;
var _elm_lang$core$Char$toLocaleUpper = _elm_lang$core$Native_Char.toLocaleUpper;
var _elm_lang$core$Char$toLower = _elm_lang$core$Native_Char.toLower;
var _elm_lang$core$Char$toUpper = _elm_lang$core$Native_Char.toUpper;
var _elm_lang$core$Char$isBetween = F3(
	function (low, high, $char) {
		var code = _elm_lang$core$Char$toCode($char);
		return (_elm_lang$core$Native_Utils.cmp(
			code,
			_elm_lang$core$Char$toCode(low)) > -1) && (_elm_lang$core$Native_Utils.cmp(
			code,
			_elm_lang$core$Char$toCode(high)) < 1);
	});
var _elm_lang$core$Char$isUpper = A2(
	_elm_lang$core$Char$isBetween,
	_elm_lang$core$Native_Utils.chr('A'),
	_elm_lang$core$Native_Utils.chr('Z'));
var _elm_lang$core$Char$isLower = A2(
	_elm_lang$core$Char$isBetween,
	_elm_lang$core$Native_Utils.chr('a'),
	_elm_lang$core$Native_Utils.chr('z'));
var _elm_lang$core$Char$isDigit = A2(
	_elm_lang$core$Char$isBetween,
	_elm_lang$core$Native_Utils.chr('0'),
	_elm_lang$core$Native_Utils.chr('9'));
var _elm_lang$core$Char$isOctDigit = A2(
	_elm_lang$core$Char$isBetween,
	_elm_lang$core$Native_Utils.chr('0'),
	_elm_lang$core$Native_Utils.chr('7'));
var _elm_lang$core$Char$isHexDigit = function ($char) {
	return _elm_lang$core$Char$isDigit($char) || (A3(
		_elm_lang$core$Char$isBetween,
		_elm_lang$core$Native_Utils.chr('a'),
		_elm_lang$core$Native_Utils.chr('f'),
		$char) || A3(
		_elm_lang$core$Char$isBetween,
		_elm_lang$core$Native_Utils.chr('A'),
		_elm_lang$core$Native_Utils.chr('F'),
		$char));
};

var _elm_lang$core$String$fromList = _elm_lang$core$Native_String.fromList;
var _elm_lang$core$String$toList = _elm_lang$core$Native_String.toList;
var _elm_lang$core$String$toFloat = _elm_lang$core$Native_String.toFloat;
var _elm_lang$core$String$toInt = _elm_lang$core$Native_String.toInt;
var _elm_lang$core$String$indices = _elm_lang$core$Native_String.indexes;
var _elm_lang$core$String$indexes = _elm_lang$core$Native_String.indexes;
var _elm_lang$core$String$endsWith = _elm_lang$core$Native_String.endsWith;
var _elm_lang$core$String$startsWith = _elm_lang$core$Native_String.startsWith;
var _elm_lang$core$String$contains = _elm_lang$core$Native_String.contains;
var _elm_lang$core$String$all = _elm_lang$core$Native_String.all;
var _elm_lang$core$String$any = _elm_lang$core$Native_String.any;
var _elm_lang$core$String$toLower = _elm_lang$core$Native_String.toLower;
var _elm_lang$core$String$toUpper = _elm_lang$core$Native_String.toUpper;
var _elm_lang$core$String$lines = _elm_lang$core$Native_String.lines;
var _elm_lang$core$String$words = _elm_lang$core$Native_String.words;
var _elm_lang$core$String$trimRight = _elm_lang$core$Native_String.trimRight;
var _elm_lang$core$String$trimLeft = _elm_lang$core$Native_String.trimLeft;
var _elm_lang$core$String$trim = _elm_lang$core$Native_String.trim;
var _elm_lang$core$String$padRight = _elm_lang$core$Native_String.padRight;
var _elm_lang$core$String$padLeft = _elm_lang$core$Native_String.padLeft;
var _elm_lang$core$String$pad = _elm_lang$core$Native_String.pad;
var _elm_lang$core$String$dropRight = _elm_lang$core$Native_String.dropRight;
var _elm_lang$core$String$dropLeft = _elm_lang$core$Native_String.dropLeft;
var _elm_lang$core$String$right = _elm_lang$core$Native_String.right;
var _elm_lang$core$String$left = _elm_lang$core$Native_String.left;
var _elm_lang$core$String$slice = _elm_lang$core$Native_String.slice;
var _elm_lang$core$String$repeat = _elm_lang$core$Native_String.repeat;
var _elm_lang$core$String$join = _elm_lang$core$Native_String.join;
var _elm_lang$core$String$split = _elm_lang$core$Native_String.split;
var _elm_lang$core$String$foldr = _elm_lang$core$Native_String.foldr;
var _elm_lang$core$String$foldl = _elm_lang$core$Native_String.foldl;
var _elm_lang$core$String$reverse = _elm_lang$core$Native_String.reverse;
var _elm_lang$core$String$filter = _elm_lang$core$Native_String.filter;
var _elm_lang$core$String$map = _elm_lang$core$Native_String.map;
var _elm_lang$core$String$length = _elm_lang$core$Native_String.length;
var _elm_lang$core$String$concat = _elm_lang$core$Native_String.concat;
var _elm_lang$core$String$append = _elm_lang$core$Native_String.append;
var _elm_lang$core$String$uncons = _elm_lang$core$Native_String.uncons;
var _elm_lang$core$String$cons = _elm_lang$core$Native_String.cons;
var _elm_lang$core$String$fromChar = function ($char) {
	return A2(_elm_lang$core$String$cons, $char, '');
};
var _elm_lang$core$String$isEmpty = _elm_lang$core$Native_String.isEmpty;

var _elm_lang$core$Dict$foldr = F3(
	function (f, acc, t) {
		foldr:
		while (true) {
			var _p0 = t;
			if (_p0.ctor === 'RBEmpty_elm_builtin') {
				return acc;
			} else {
				var _v1 = f,
					_v2 = A3(
					f,
					_p0._1,
					_p0._2,
					A3(_elm_lang$core$Dict$foldr, f, acc, _p0._4)),
					_v3 = _p0._3;
				f = _v1;
				acc = _v2;
				t = _v3;
				continue foldr;
			}
		}
	});
var _elm_lang$core$Dict$keys = function (dict) {
	return A3(
		_elm_lang$core$Dict$foldr,
		F3(
			function (key, value, keyList) {
				return A2(_elm_lang$core$List_ops['::'], key, keyList);
			}),
		_elm_lang$core$Native_List.fromArray(
			[]),
		dict);
};
var _elm_lang$core$Dict$values = function (dict) {
	return A3(
		_elm_lang$core$Dict$foldr,
		F3(
			function (key, value, valueList) {
				return A2(_elm_lang$core$List_ops['::'], value, valueList);
			}),
		_elm_lang$core$Native_List.fromArray(
			[]),
		dict);
};
var _elm_lang$core$Dict$toList = function (dict) {
	return A3(
		_elm_lang$core$Dict$foldr,
		F3(
			function (key, value, list) {
				return A2(
					_elm_lang$core$List_ops['::'],
					{ctor: '_Tuple2', _0: key, _1: value},
					list);
			}),
		_elm_lang$core$Native_List.fromArray(
			[]),
		dict);
};
var _elm_lang$core$Dict$foldl = F3(
	function (f, acc, dict) {
		foldl:
		while (true) {
			var _p1 = dict;
			if (_p1.ctor === 'RBEmpty_elm_builtin') {
				return acc;
			} else {
				var _v5 = f,
					_v6 = A3(
					f,
					_p1._1,
					_p1._2,
					A3(_elm_lang$core$Dict$foldl, f, acc, _p1._3)),
					_v7 = _p1._4;
				f = _v5;
				acc = _v6;
				dict = _v7;
				continue foldl;
			}
		}
	});
var _elm_lang$core$Dict$merge = F6(
	function (leftStep, bothStep, rightStep, leftDict, rightDict, initialResult) {
		var stepState = F3(
			function (rKey, rValue, _p2) {
				stepState:
				while (true) {
					var _p3 = _p2;
					var _p9 = _p3._1;
					var _p8 = _p3._0;
					var _p4 = _p8;
					if (_p4.ctor === '[]') {
						return {
							ctor: '_Tuple2',
							_0: _p8,
							_1: A3(rightStep, rKey, rValue, _p9)
						};
					} else {
						var _p7 = _p4._1;
						var _p6 = _p4._0._1;
						var _p5 = _p4._0._0;
						if (_elm_lang$core$Native_Utils.cmp(_p5, rKey) < 0) {
							var _v10 = rKey,
								_v11 = rValue,
								_v12 = {
								ctor: '_Tuple2',
								_0: _p7,
								_1: A3(leftStep, _p5, _p6, _p9)
							};
							rKey = _v10;
							rValue = _v11;
							_p2 = _v12;
							continue stepState;
						} else {
							if (_elm_lang$core$Native_Utils.cmp(_p5, rKey) > 0) {
								return {
									ctor: '_Tuple2',
									_0: _p8,
									_1: A3(rightStep, rKey, rValue, _p9)
								};
							} else {
								return {
									ctor: '_Tuple2',
									_0: _p7,
									_1: A4(bothStep, _p5, _p6, rValue, _p9)
								};
							}
						}
					}
				}
			});
		var _p10 = A3(
			_elm_lang$core$Dict$foldl,
			stepState,
			{
				ctor: '_Tuple2',
				_0: _elm_lang$core$Dict$toList(leftDict),
				_1: initialResult
			},
			rightDict);
		var leftovers = _p10._0;
		var intermediateResult = _p10._1;
		return A3(
			_elm_lang$core$List$foldl,
			F2(
				function (_p11, result) {
					var _p12 = _p11;
					return A3(leftStep, _p12._0, _p12._1, result);
				}),
			intermediateResult,
			leftovers);
	});
var _elm_lang$core$Dict$reportRemBug = F4(
	function (msg, c, lgot, rgot) {
		return _elm_lang$core$Native_Debug.crash(
			_elm_lang$core$String$concat(
				_elm_lang$core$Native_List.fromArray(
					[
						'Internal red-black tree invariant violated, expected ',
						msg,
						' and got ',
						_elm_lang$core$Basics$toString(c),
						'/',
						lgot,
						'/',
						rgot,
						'\nPlease report this bug to <https://github.com/elm-lang/core/issues>'
					])));
	});
var _elm_lang$core$Dict$isBBlack = function (dict) {
	var _p13 = dict;
	_v14_2:
	do {
		if (_p13.ctor === 'RBNode_elm_builtin') {
			if (_p13._0.ctor === 'BBlack') {
				return true;
			} else {
				break _v14_2;
			}
		} else {
			if (_p13._0.ctor === 'LBBlack') {
				return true;
			} else {
				break _v14_2;
			}
		}
	} while(false);
	return false;
};
var _elm_lang$core$Dict$sizeHelp = F2(
	function (n, dict) {
		sizeHelp:
		while (true) {
			var _p14 = dict;
			if (_p14.ctor === 'RBEmpty_elm_builtin') {
				return n;
			} else {
				var _v16 = A2(_elm_lang$core$Dict$sizeHelp, n + 1, _p14._4),
					_v17 = _p14._3;
				n = _v16;
				dict = _v17;
				continue sizeHelp;
			}
		}
	});
var _elm_lang$core$Dict$size = function (dict) {
	return A2(_elm_lang$core$Dict$sizeHelp, 0, dict);
};
var _elm_lang$core$Dict$get = F2(
	function (targetKey, dict) {
		get:
		while (true) {
			var _p15 = dict;
			if (_p15.ctor === 'RBEmpty_elm_builtin') {
				return _elm_lang$core$Maybe$Nothing;
			} else {
				var _p16 = A2(_elm_lang$core$Basics$compare, targetKey, _p15._1);
				switch (_p16.ctor) {
					case 'LT':
						var _v20 = targetKey,
							_v21 = _p15._3;
						targetKey = _v20;
						dict = _v21;
						continue get;
					case 'EQ':
						return _elm_lang$core$Maybe$Just(_p15._2);
					default:
						var _v22 = targetKey,
							_v23 = _p15._4;
						targetKey = _v22;
						dict = _v23;
						continue get;
				}
			}
		}
	});
var _elm_lang$core$Dict$member = F2(
	function (key, dict) {
		var _p17 = A2(_elm_lang$core$Dict$get, key, dict);
		if (_p17.ctor === 'Just') {
			return true;
		} else {
			return false;
		}
	});
var _elm_lang$core$Dict$maxWithDefault = F3(
	function (k, v, r) {
		maxWithDefault:
		while (true) {
			var _p18 = r;
			if (_p18.ctor === 'RBEmpty_elm_builtin') {
				return {ctor: '_Tuple2', _0: k, _1: v};
			} else {
				var _v26 = _p18._1,
					_v27 = _p18._2,
					_v28 = _p18._4;
				k = _v26;
				v = _v27;
				r = _v28;
				continue maxWithDefault;
			}
		}
	});
var _elm_lang$core$Dict$NBlack = {ctor: 'NBlack'};
var _elm_lang$core$Dict$BBlack = {ctor: 'BBlack'};
var _elm_lang$core$Dict$Black = {ctor: 'Black'};
var _elm_lang$core$Dict$blackish = function (t) {
	var _p19 = t;
	if (_p19.ctor === 'RBNode_elm_builtin') {
		var _p20 = _p19._0;
		return _elm_lang$core$Native_Utils.eq(_p20, _elm_lang$core$Dict$Black) || _elm_lang$core$Native_Utils.eq(_p20, _elm_lang$core$Dict$BBlack);
	} else {
		return true;
	}
};
var _elm_lang$core$Dict$Red = {ctor: 'Red'};
var _elm_lang$core$Dict$moreBlack = function (color) {
	var _p21 = color;
	switch (_p21.ctor) {
		case 'Black':
			return _elm_lang$core$Dict$BBlack;
		case 'Red':
			return _elm_lang$core$Dict$Black;
		case 'NBlack':
			return _elm_lang$core$Dict$Red;
		default:
			return _elm_lang$core$Native_Debug.crash('Can\'t make a double black node more black!');
	}
};
var _elm_lang$core$Dict$lessBlack = function (color) {
	var _p22 = color;
	switch (_p22.ctor) {
		case 'BBlack':
			return _elm_lang$core$Dict$Black;
		case 'Black':
			return _elm_lang$core$Dict$Red;
		case 'Red':
			return _elm_lang$core$Dict$NBlack;
		default:
			return _elm_lang$core$Native_Debug.crash('Can\'t make a negative black node less black!');
	}
};
var _elm_lang$core$Dict$LBBlack = {ctor: 'LBBlack'};
var _elm_lang$core$Dict$LBlack = {ctor: 'LBlack'};
var _elm_lang$core$Dict$RBEmpty_elm_builtin = function (a) {
	return {ctor: 'RBEmpty_elm_builtin', _0: a};
};
var _elm_lang$core$Dict$empty = _elm_lang$core$Dict$RBEmpty_elm_builtin(_elm_lang$core$Dict$LBlack);
var _elm_lang$core$Dict$isEmpty = function (dict) {
	return _elm_lang$core$Native_Utils.eq(dict, _elm_lang$core$Dict$empty);
};
var _elm_lang$core$Dict$RBNode_elm_builtin = F5(
	function (a, b, c, d, e) {
		return {ctor: 'RBNode_elm_builtin', _0: a, _1: b, _2: c, _3: d, _4: e};
	});
var _elm_lang$core$Dict$ensureBlackRoot = function (dict) {
	var _p23 = dict;
	if ((_p23.ctor === 'RBNode_elm_builtin') && (_p23._0.ctor === 'Red')) {
		return A5(_elm_lang$core$Dict$RBNode_elm_builtin, _elm_lang$core$Dict$Black, _p23._1, _p23._2, _p23._3, _p23._4);
	} else {
		return dict;
	}
};
var _elm_lang$core$Dict$lessBlackTree = function (dict) {
	var _p24 = dict;
	if (_p24.ctor === 'RBNode_elm_builtin') {
		return A5(
			_elm_lang$core$Dict$RBNode_elm_builtin,
			_elm_lang$core$Dict$lessBlack(_p24._0),
			_p24._1,
			_p24._2,
			_p24._3,
			_p24._4);
	} else {
		return _elm_lang$core$Dict$RBEmpty_elm_builtin(_elm_lang$core$Dict$LBlack);
	}
};
var _elm_lang$core$Dict$balancedTree = function (col) {
	return function (xk) {
		return function (xv) {
			return function (yk) {
				return function (yv) {
					return function (zk) {
						return function (zv) {
							return function (a) {
								return function (b) {
									return function (c) {
										return function (d) {
											return A5(
												_elm_lang$core$Dict$RBNode_elm_builtin,
												_elm_lang$core$Dict$lessBlack(col),
												yk,
												yv,
												A5(_elm_lang$core$Dict$RBNode_elm_builtin, _elm_lang$core$Dict$Black, xk, xv, a, b),
												A5(_elm_lang$core$Dict$RBNode_elm_builtin, _elm_lang$core$Dict$Black, zk, zv, c, d));
										};
									};
								};
							};
						};
					};
				};
			};
		};
	};
};
var _elm_lang$core$Dict$blacken = function (t) {
	var _p25 = t;
	if (_p25.ctor === 'RBEmpty_elm_builtin') {
		return _elm_lang$core$Dict$RBEmpty_elm_builtin(_elm_lang$core$Dict$LBlack);
	} else {
		return A5(_elm_lang$core$Dict$RBNode_elm_builtin, _elm_lang$core$Dict$Black, _p25._1, _p25._2, _p25._3, _p25._4);
	}
};
var _elm_lang$core$Dict$redden = function (t) {
	var _p26 = t;
	if (_p26.ctor === 'RBEmpty_elm_builtin') {
		return _elm_lang$core$Native_Debug.crash('can\'t make a Leaf red');
	} else {
		return A5(_elm_lang$core$Dict$RBNode_elm_builtin, _elm_lang$core$Dict$Red, _p26._1, _p26._2, _p26._3, _p26._4);
	}
};
var _elm_lang$core$Dict$balanceHelp = function (tree) {
	var _p27 = tree;
	_v36_6:
	do {
		_v36_5:
		do {
			_v36_4:
			do {
				_v36_3:
				do {
					_v36_2:
					do {
						_v36_1:
						do {
							_v36_0:
							do {
								if (_p27.ctor === 'RBNode_elm_builtin') {
									if (_p27._3.ctor === 'RBNode_elm_builtin') {
										if (_p27._4.ctor === 'RBNode_elm_builtin') {
											switch (_p27._3._0.ctor) {
												case 'Red':
													switch (_p27._4._0.ctor) {
														case 'Red':
															if ((_p27._3._3.ctor === 'RBNode_elm_builtin') && (_p27._3._3._0.ctor === 'Red')) {
																break _v36_0;
															} else {
																if ((_p27._3._4.ctor === 'RBNode_elm_builtin') && (_p27._3._4._0.ctor === 'Red')) {
																	break _v36_1;
																} else {
																	if ((_p27._4._3.ctor === 'RBNode_elm_builtin') && (_p27._4._3._0.ctor === 'Red')) {
																		break _v36_2;
																	} else {
																		if ((_p27._4._4.ctor === 'RBNode_elm_builtin') && (_p27._4._4._0.ctor === 'Red')) {
																			break _v36_3;
																		} else {
																			break _v36_6;
																		}
																	}
																}
															}
														case 'NBlack':
															if ((_p27._3._3.ctor === 'RBNode_elm_builtin') && (_p27._3._3._0.ctor === 'Red')) {
																break _v36_0;
															} else {
																if ((_p27._3._4.ctor === 'RBNode_elm_builtin') && (_p27._3._4._0.ctor === 'Red')) {
																	break _v36_1;
																} else {
																	if (((((_p27._0.ctor === 'BBlack') && (_p27._4._3.ctor === 'RBNode_elm_builtin')) && (_p27._4._3._0.ctor === 'Black')) && (_p27._4._4.ctor === 'RBNode_elm_builtin')) && (_p27._4._4._0.ctor === 'Black')) {
																		break _v36_4;
																	} else {
																		break _v36_6;
																	}
																}
															}
														default:
															if ((_p27._3._3.ctor === 'RBNode_elm_builtin') && (_p27._3._3._0.ctor === 'Red')) {
																break _v36_0;
															} else {
																if ((_p27._3._4.ctor === 'RBNode_elm_builtin') && (_p27._3._4._0.ctor === 'Red')) {
																	break _v36_1;
																} else {
																	break _v36_6;
																}
															}
													}
												case 'NBlack':
													switch (_p27._4._0.ctor) {
														case 'Red':
															if ((_p27._4._3.ctor === 'RBNode_elm_builtin') && (_p27._4._3._0.ctor === 'Red')) {
																break _v36_2;
															} else {
																if ((_p27._4._4.ctor === 'RBNode_elm_builtin') && (_p27._4._4._0.ctor === 'Red')) {
																	break _v36_3;
																} else {
																	if (((((_p27._0.ctor === 'BBlack') && (_p27._3._3.ctor === 'RBNode_elm_builtin')) && (_p27._3._3._0.ctor === 'Black')) && (_p27._3._4.ctor === 'RBNode_elm_builtin')) && (_p27._3._4._0.ctor === 'Black')) {
																		break _v36_5;
																	} else {
																		break _v36_6;
																	}
																}
															}
														case 'NBlack':
															if (_p27._0.ctor === 'BBlack') {
																if ((((_p27._4._3.ctor === 'RBNode_elm_builtin') && (_p27._4._3._0.ctor === 'Black')) && (_p27._4._4.ctor === 'RBNode_elm_builtin')) && (_p27._4._4._0.ctor === 'Black')) {
																	break _v36_4;
																} else {
																	if ((((_p27._3._3.ctor === 'RBNode_elm_builtin') && (_p27._3._3._0.ctor === 'Black')) && (_p27._3._4.ctor === 'RBNode_elm_builtin')) && (_p27._3._4._0.ctor === 'Black')) {
																		break _v36_5;
																	} else {
																		break _v36_6;
																	}
																}
															} else {
																break _v36_6;
															}
														default:
															if (((((_p27._0.ctor === 'BBlack') && (_p27._3._3.ctor === 'RBNode_elm_builtin')) && (_p27._3._3._0.ctor === 'Black')) && (_p27._3._4.ctor === 'RBNode_elm_builtin')) && (_p27._3._4._0.ctor === 'Black')) {
																break _v36_5;
															} else {
																break _v36_6;
															}
													}
												default:
													switch (_p27._4._0.ctor) {
														case 'Red':
															if ((_p27._4._3.ctor === 'RBNode_elm_builtin') && (_p27._4._3._0.ctor === 'Red')) {
																break _v36_2;
															} else {
																if ((_p27._4._4.ctor === 'RBNode_elm_builtin') && (_p27._4._4._0.ctor === 'Red')) {
																	break _v36_3;
																} else {
																	break _v36_6;
																}
															}
														case 'NBlack':
															if (((((_p27._0.ctor === 'BBlack') && (_p27._4._3.ctor === 'RBNode_elm_builtin')) && (_p27._4._3._0.ctor === 'Black')) && (_p27._4._4.ctor === 'RBNode_elm_builtin')) && (_p27._4._4._0.ctor === 'Black')) {
																break _v36_4;
															} else {
																break _v36_6;
															}
														default:
															break _v36_6;
													}
											}
										} else {
											switch (_p27._3._0.ctor) {
												case 'Red':
													if ((_p27._3._3.ctor === 'RBNode_elm_builtin') && (_p27._3._3._0.ctor === 'Red')) {
														break _v36_0;
													} else {
														if ((_p27._3._4.ctor === 'RBNode_elm_builtin') && (_p27._3._4._0.ctor === 'Red')) {
															break _v36_1;
														} else {
															break _v36_6;
														}
													}
												case 'NBlack':
													if (((((_p27._0.ctor === 'BBlack') && (_p27._3._3.ctor === 'RBNode_elm_builtin')) && (_p27._3._3._0.ctor === 'Black')) && (_p27._3._4.ctor === 'RBNode_elm_builtin')) && (_p27._3._4._0.ctor === 'Black')) {
														break _v36_5;
													} else {
														break _v36_6;
													}
												default:
													break _v36_6;
											}
										}
									} else {
										if (_p27._4.ctor === 'RBNode_elm_builtin') {
											switch (_p27._4._0.ctor) {
												case 'Red':
													if ((_p27._4._3.ctor === 'RBNode_elm_builtin') && (_p27._4._3._0.ctor === 'Red')) {
														break _v36_2;
													} else {
														if ((_p27._4._4.ctor === 'RBNode_elm_builtin') && (_p27._4._4._0.ctor === 'Red')) {
															break _v36_3;
														} else {
															break _v36_6;
														}
													}
												case 'NBlack':
													if (((((_p27._0.ctor === 'BBlack') && (_p27._4._3.ctor === 'RBNode_elm_builtin')) && (_p27._4._3._0.ctor === 'Black')) && (_p27._4._4.ctor === 'RBNode_elm_builtin')) && (_p27._4._4._0.ctor === 'Black')) {
														break _v36_4;
													} else {
														break _v36_6;
													}
												default:
													break _v36_6;
											}
										} else {
											break _v36_6;
										}
									}
								} else {
									break _v36_6;
								}
							} while(false);
							return _elm_lang$core$Dict$balancedTree(_p27._0)(_p27._3._3._1)(_p27._3._3._2)(_p27._3._1)(_p27._3._2)(_p27._1)(_p27._2)(_p27._3._3._3)(_p27._3._3._4)(_p27._3._4)(_p27._4);
						} while(false);
						return _elm_lang$core$Dict$balancedTree(_p27._0)(_p27._3._1)(_p27._3._2)(_p27._3._4._1)(_p27._3._4._2)(_p27._1)(_p27._2)(_p27._3._3)(_p27._3._4._3)(_p27._3._4._4)(_p27._4);
					} while(false);
					return _elm_lang$core$Dict$balancedTree(_p27._0)(_p27._1)(_p27._2)(_p27._4._3._1)(_p27._4._3._2)(_p27._4._1)(_p27._4._2)(_p27._3)(_p27._4._3._3)(_p27._4._3._4)(_p27._4._4);
				} while(false);
				return _elm_lang$core$Dict$balancedTree(_p27._0)(_p27._1)(_p27._2)(_p27._4._1)(_p27._4._2)(_p27._4._4._1)(_p27._4._4._2)(_p27._3)(_p27._4._3)(_p27._4._4._3)(_p27._4._4._4);
			} while(false);
			return A5(
				_elm_lang$core$Dict$RBNode_elm_builtin,
				_elm_lang$core$Dict$Black,
				_p27._4._3._1,
				_p27._4._3._2,
				A5(_elm_lang$core$Dict$RBNode_elm_builtin, _elm_lang$core$Dict$Black, _p27._1, _p27._2, _p27._3, _p27._4._3._3),
				A5(
					_elm_lang$core$Dict$balance,
					_elm_lang$core$Dict$Black,
					_p27._4._1,
					_p27._4._2,
					_p27._4._3._4,
					_elm_lang$core$Dict$redden(_p27._4._4)));
		} while(false);
		return A5(
			_elm_lang$core$Dict$RBNode_elm_builtin,
			_elm_lang$core$Dict$Black,
			_p27._3._4._1,
			_p27._3._4._2,
			A5(
				_elm_lang$core$Dict$balance,
				_elm_lang$core$Dict$Black,
				_p27._3._1,
				_p27._3._2,
				_elm_lang$core$Dict$redden(_p27._3._3),
				_p27._3._4._3),
			A5(_elm_lang$core$Dict$RBNode_elm_builtin, _elm_lang$core$Dict$Black, _p27._1, _p27._2, _p27._3._4._4, _p27._4));
	} while(false);
	return tree;
};
var _elm_lang$core$Dict$balance = F5(
	function (c, k, v, l, r) {
		var tree = A5(_elm_lang$core$Dict$RBNode_elm_builtin, c, k, v, l, r);
		return _elm_lang$core$Dict$blackish(tree) ? _elm_lang$core$Dict$balanceHelp(tree) : tree;
	});
var _elm_lang$core$Dict$bubble = F5(
	function (c, k, v, l, r) {
		return (_elm_lang$core$Dict$isBBlack(l) || _elm_lang$core$Dict$isBBlack(r)) ? A5(
			_elm_lang$core$Dict$balance,
			_elm_lang$core$Dict$moreBlack(c),
			k,
			v,
			_elm_lang$core$Dict$lessBlackTree(l),
			_elm_lang$core$Dict$lessBlackTree(r)) : A5(_elm_lang$core$Dict$RBNode_elm_builtin, c, k, v, l, r);
	});
var _elm_lang$core$Dict$removeMax = F5(
	function (c, k, v, l, r) {
		var _p28 = r;
		if (_p28.ctor === 'RBEmpty_elm_builtin') {
			return A3(_elm_lang$core$Dict$rem, c, l, r);
		} else {
			return A5(
				_elm_lang$core$Dict$bubble,
				c,
				k,
				v,
				l,
				A5(_elm_lang$core$Dict$removeMax, _p28._0, _p28._1, _p28._2, _p28._3, _p28._4));
		}
	});
var _elm_lang$core$Dict$rem = F3(
	function (c, l, r) {
		var _p29 = {ctor: '_Tuple2', _0: l, _1: r};
		if (_p29._0.ctor === 'RBEmpty_elm_builtin') {
			if (_p29._1.ctor === 'RBEmpty_elm_builtin') {
				var _p30 = c;
				switch (_p30.ctor) {
					case 'Red':
						return _elm_lang$core$Dict$RBEmpty_elm_builtin(_elm_lang$core$Dict$LBlack);
					case 'Black':
						return _elm_lang$core$Dict$RBEmpty_elm_builtin(_elm_lang$core$Dict$LBBlack);
					default:
						return _elm_lang$core$Native_Debug.crash('cannot have bblack or nblack nodes at this point');
				}
			} else {
				var _p33 = _p29._1._0;
				var _p32 = _p29._0._0;
				var _p31 = {ctor: '_Tuple3', _0: c, _1: _p32, _2: _p33};
				if ((((_p31.ctor === '_Tuple3') && (_p31._0.ctor === 'Black')) && (_p31._1.ctor === 'LBlack')) && (_p31._2.ctor === 'Red')) {
					return A5(_elm_lang$core$Dict$RBNode_elm_builtin, _elm_lang$core$Dict$Black, _p29._1._1, _p29._1._2, _p29._1._3, _p29._1._4);
				} else {
					return A4(
						_elm_lang$core$Dict$reportRemBug,
						'Black/LBlack/Red',
						c,
						_elm_lang$core$Basics$toString(_p32),
						_elm_lang$core$Basics$toString(_p33));
				}
			}
		} else {
			if (_p29._1.ctor === 'RBEmpty_elm_builtin') {
				var _p36 = _p29._1._0;
				var _p35 = _p29._0._0;
				var _p34 = {ctor: '_Tuple3', _0: c, _1: _p35, _2: _p36};
				if ((((_p34.ctor === '_Tuple3') && (_p34._0.ctor === 'Black')) && (_p34._1.ctor === 'Red')) && (_p34._2.ctor === 'LBlack')) {
					return A5(_elm_lang$core$Dict$RBNode_elm_builtin, _elm_lang$core$Dict$Black, _p29._0._1, _p29._0._2, _p29._0._3, _p29._0._4);
				} else {
					return A4(
						_elm_lang$core$Dict$reportRemBug,
						'Black/Red/LBlack',
						c,
						_elm_lang$core$Basics$toString(_p35),
						_elm_lang$core$Basics$toString(_p36));
				}
			} else {
				var _p40 = _p29._0._2;
				var _p39 = _p29._0._4;
				var _p38 = _p29._0._1;
				var l$ = A5(_elm_lang$core$Dict$removeMax, _p29._0._0, _p38, _p40, _p29._0._3, _p39);
				var _p37 = A3(_elm_lang$core$Dict$maxWithDefault, _p38, _p40, _p39);
				var k = _p37._0;
				var v = _p37._1;
				return A5(_elm_lang$core$Dict$bubble, c, k, v, l$, r);
			}
		}
	});
var _elm_lang$core$Dict$map = F2(
	function (f, dict) {
		var _p41 = dict;
		if (_p41.ctor === 'RBEmpty_elm_builtin') {
			return _elm_lang$core$Dict$RBEmpty_elm_builtin(_elm_lang$core$Dict$LBlack);
		} else {
			var _p42 = _p41._1;
			return A5(
				_elm_lang$core$Dict$RBNode_elm_builtin,
				_p41._0,
				_p42,
				A2(f, _p42, _p41._2),
				A2(_elm_lang$core$Dict$map, f, _p41._3),
				A2(_elm_lang$core$Dict$map, f, _p41._4));
		}
	});
var _elm_lang$core$Dict$Same = {ctor: 'Same'};
var _elm_lang$core$Dict$Remove = {ctor: 'Remove'};
var _elm_lang$core$Dict$Insert = {ctor: 'Insert'};
var _elm_lang$core$Dict$update = F3(
	function (k, alter, dict) {
		var up = function (dict) {
			var _p43 = dict;
			if (_p43.ctor === 'RBEmpty_elm_builtin') {
				var _p44 = alter(_elm_lang$core$Maybe$Nothing);
				if (_p44.ctor === 'Nothing') {
					return {ctor: '_Tuple2', _0: _elm_lang$core$Dict$Same, _1: _elm_lang$core$Dict$empty};
				} else {
					return {
						ctor: '_Tuple2',
						_0: _elm_lang$core$Dict$Insert,
						_1: A5(_elm_lang$core$Dict$RBNode_elm_builtin, _elm_lang$core$Dict$Red, k, _p44._0, _elm_lang$core$Dict$empty, _elm_lang$core$Dict$empty)
					};
				}
			} else {
				var _p55 = _p43._2;
				var _p54 = _p43._4;
				var _p53 = _p43._3;
				var _p52 = _p43._1;
				var _p51 = _p43._0;
				var _p45 = A2(_elm_lang$core$Basics$compare, k, _p52);
				switch (_p45.ctor) {
					case 'EQ':
						var _p46 = alter(
							_elm_lang$core$Maybe$Just(_p55));
						if (_p46.ctor === 'Nothing') {
							return {
								ctor: '_Tuple2',
								_0: _elm_lang$core$Dict$Remove,
								_1: A3(_elm_lang$core$Dict$rem, _p51, _p53, _p54)
							};
						} else {
							return {
								ctor: '_Tuple2',
								_0: _elm_lang$core$Dict$Same,
								_1: A5(_elm_lang$core$Dict$RBNode_elm_builtin, _p51, _p52, _p46._0, _p53, _p54)
							};
						}
					case 'LT':
						var _p47 = up(_p53);
						var flag = _p47._0;
						var newLeft = _p47._1;
						var _p48 = flag;
						switch (_p48.ctor) {
							case 'Same':
								return {
									ctor: '_Tuple2',
									_0: _elm_lang$core$Dict$Same,
									_1: A5(_elm_lang$core$Dict$RBNode_elm_builtin, _p51, _p52, _p55, newLeft, _p54)
								};
							case 'Insert':
								return {
									ctor: '_Tuple2',
									_0: _elm_lang$core$Dict$Insert,
									_1: A5(_elm_lang$core$Dict$balance, _p51, _p52, _p55, newLeft, _p54)
								};
							default:
								return {
									ctor: '_Tuple2',
									_0: _elm_lang$core$Dict$Remove,
									_1: A5(_elm_lang$core$Dict$bubble, _p51, _p52, _p55, newLeft, _p54)
								};
						}
					default:
						var _p49 = up(_p54);
						var flag = _p49._0;
						var newRight = _p49._1;
						var _p50 = flag;
						switch (_p50.ctor) {
							case 'Same':
								return {
									ctor: '_Tuple2',
									_0: _elm_lang$core$Dict$Same,
									_1: A5(_elm_lang$core$Dict$RBNode_elm_builtin, _p51, _p52, _p55, _p53, newRight)
								};
							case 'Insert':
								return {
									ctor: '_Tuple2',
									_0: _elm_lang$core$Dict$Insert,
									_1: A5(_elm_lang$core$Dict$balance, _p51, _p52, _p55, _p53, newRight)
								};
							default:
								return {
									ctor: '_Tuple2',
									_0: _elm_lang$core$Dict$Remove,
									_1: A5(_elm_lang$core$Dict$bubble, _p51, _p52, _p55, _p53, newRight)
								};
						}
				}
			}
		};
		var _p56 = up(dict);
		var flag = _p56._0;
		var updatedDict = _p56._1;
		var _p57 = flag;
		switch (_p57.ctor) {
			case 'Same':
				return updatedDict;
			case 'Insert':
				return _elm_lang$core$Dict$ensureBlackRoot(updatedDict);
			default:
				return _elm_lang$core$Dict$blacken(updatedDict);
		}
	});
var _elm_lang$core$Dict$insert = F3(
	function (key, value, dict) {
		return A3(
			_elm_lang$core$Dict$update,
			key,
			_elm_lang$core$Basics$always(
				_elm_lang$core$Maybe$Just(value)),
			dict);
	});
var _elm_lang$core$Dict$singleton = F2(
	function (key, value) {
		return A3(_elm_lang$core$Dict$insert, key, value, _elm_lang$core$Dict$empty);
	});
var _elm_lang$core$Dict$union = F2(
	function (t1, t2) {
		return A3(_elm_lang$core$Dict$foldl, _elm_lang$core$Dict$insert, t2, t1);
	});
var _elm_lang$core$Dict$filter = F2(
	function (predicate, dictionary) {
		var add = F3(
			function (key, value, dict) {
				return A2(predicate, key, value) ? A3(_elm_lang$core$Dict$insert, key, value, dict) : dict;
			});
		return A3(_elm_lang$core$Dict$foldl, add, _elm_lang$core$Dict$empty, dictionary);
	});
var _elm_lang$core$Dict$intersect = F2(
	function (t1, t2) {
		return A2(
			_elm_lang$core$Dict$filter,
			F2(
				function (k, _p58) {
					return A2(_elm_lang$core$Dict$member, k, t2);
				}),
			t1);
	});
var _elm_lang$core$Dict$partition = F2(
	function (predicate, dict) {
		var add = F3(
			function (key, value, _p59) {
				var _p60 = _p59;
				var _p62 = _p60._1;
				var _p61 = _p60._0;
				return A2(predicate, key, value) ? {
					ctor: '_Tuple2',
					_0: A3(_elm_lang$core$Dict$insert, key, value, _p61),
					_1: _p62
				} : {
					ctor: '_Tuple2',
					_0: _p61,
					_1: A3(_elm_lang$core$Dict$insert, key, value, _p62)
				};
			});
		return A3(
			_elm_lang$core$Dict$foldl,
			add,
			{ctor: '_Tuple2', _0: _elm_lang$core$Dict$empty, _1: _elm_lang$core$Dict$empty},
			dict);
	});
var _elm_lang$core$Dict$fromList = function (assocs) {
	return A3(
		_elm_lang$core$List$foldl,
		F2(
			function (_p63, dict) {
				var _p64 = _p63;
				return A3(_elm_lang$core$Dict$insert, _p64._0, _p64._1, dict);
			}),
		_elm_lang$core$Dict$empty,
		assocs);
};
var _elm_lang$core$Dict$remove = F2(
	function (key, dict) {
		return A3(
			_elm_lang$core$Dict$update,
			key,
			_elm_lang$core$Basics$always(_elm_lang$core$Maybe$Nothing),
			dict);
	});
var _elm_lang$core$Dict$diff = F2(
	function (t1, t2) {
		return A3(
			_elm_lang$core$Dict$foldl,
			F3(
				function (k, v, t) {
					return A2(_elm_lang$core$Dict$remove, k, t);
				}),
			t1,
			t2);
	});

//import Native.Scheduler //

var _elm_lang$core$Native_Time = function() {

var now = _elm_lang$core$Native_Scheduler.nativeBinding(function(callback)
{
	callback(_elm_lang$core$Native_Scheduler.succeed(Date.now()));
});

function setInterval_(interval, task)
{
	return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback)
	{
		var id = setInterval(function() {
			_elm_lang$core$Native_Scheduler.rawSpawn(task);
		}, interval);

		return function() { clearInterval(id); };
	});
}

return {
	now: now,
	setInterval_: F2(setInterval_)
};

}();
var _elm_lang$core$Time$setInterval = _elm_lang$core$Native_Time.setInterval_;
var _elm_lang$core$Time$spawnHelp = F3(
	function (router, intervals, processes) {
		var _p0 = intervals;
		if (_p0.ctor === '[]') {
			return _elm_lang$core$Task$succeed(processes);
		} else {
			var _p1 = _p0._0;
			return A2(
				_elm_lang$core$Task$andThen,
				_elm_lang$core$Native_Scheduler.spawn(
					A2(
						_elm_lang$core$Time$setInterval,
						_p1,
						A2(_elm_lang$core$Platform$sendToSelf, router, _p1))),
				function (id) {
					return A3(
						_elm_lang$core$Time$spawnHelp,
						router,
						_p0._1,
						A3(_elm_lang$core$Dict$insert, _p1, id, processes));
				});
		}
	});
var _elm_lang$core$Time$addMySub = F2(
	function (_p2, state) {
		var _p3 = _p2;
		var _p6 = _p3._1;
		var _p5 = _p3._0;
		var _p4 = A2(_elm_lang$core$Dict$get, _p5, state);
		if (_p4.ctor === 'Nothing') {
			return A3(
				_elm_lang$core$Dict$insert,
				_p5,
				_elm_lang$core$Native_List.fromArray(
					[_p6]),
				state);
		} else {
			return A3(
				_elm_lang$core$Dict$insert,
				_p5,
				A2(_elm_lang$core$List_ops['::'], _p6, _p4._0),
				state);
		}
	});
var _elm_lang$core$Time$inMilliseconds = function (t) {
	return t;
};
var _elm_lang$core$Time$millisecond = 1;
var _elm_lang$core$Time$second = 1000 * _elm_lang$core$Time$millisecond;
var _elm_lang$core$Time$minute = 60 * _elm_lang$core$Time$second;
var _elm_lang$core$Time$hour = 60 * _elm_lang$core$Time$minute;
var _elm_lang$core$Time$inHours = function (t) {
	return t / _elm_lang$core$Time$hour;
};
var _elm_lang$core$Time$inMinutes = function (t) {
	return t / _elm_lang$core$Time$minute;
};
var _elm_lang$core$Time$inSeconds = function (t) {
	return t / _elm_lang$core$Time$second;
};
var _elm_lang$core$Time$now = _elm_lang$core$Native_Time.now;
var _elm_lang$core$Time$onSelfMsg = F3(
	function (router, interval, state) {
		var _p7 = A2(_elm_lang$core$Dict$get, interval, state.taggers);
		if (_p7.ctor === 'Nothing') {
			return _elm_lang$core$Task$succeed(state);
		} else {
			return A2(
				_elm_lang$core$Task$andThen,
				_elm_lang$core$Time$now,
				function (time) {
					return A2(
						_elm_lang$core$Task$andThen,
						_elm_lang$core$Task$sequence(
							A2(
								_elm_lang$core$List$map,
								function (tagger) {
									return A2(
										_elm_lang$core$Platform$sendToApp,
										router,
										tagger(time));
								},
								_p7._0)),
						function (_p8) {
							return _elm_lang$core$Task$succeed(state);
						});
				});
		}
	});
var _elm_lang$core$Time$subscription = _elm_lang$core$Native_Platform.leaf('Time');
var _elm_lang$core$Time$State = F2(
	function (a, b) {
		return {taggers: a, processes: b};
	});
var _elm_lang$core$Time$init = _elm_lang$core$Task$succeed(
	A2(_elm_lang$core$Time$State, _elm_lang$core$Dict$empty, _elm_lang$core$Dict$empty));
var _elm_lang$core$Time$onEffects = F3(
	function (router, subs, _p9) {
		var _p10 = _p9;
		var rightStep = F3(
			function (_p12, id, _p11) {
				var _p13 = _p11;
				return {
					ctor: '_Tuple3',
					_0: _p13._0,
					_1: _p13._1,
					_2: A2(
						_elm_lang$core$Task$andThen,
						_elm_lang$core$Native_Scheduler.kill(id),
						function (_p14) {
							return _p13._2;
						})
				};
			});
		var bothStep = F4(
			function (interval, taggers, id, _p15) {
				var _p16 = _p15;
				return {
					ctor: '_Tuple3',
					_0: _p16._0,
					_1: A3(_elm_lang$core$Dict$insert, interval, id, _p16._1),
					_2: _p16._2
				};
			});
		var leftStep = F3(
			function (interval, taggers, _p17) {
				var _p18 = _p17;
				return {
					ctor: '_Tuple3',
					_0: A2(_elm_lang$core$List_ops['::'], interval, _p18._0),
					_1: _p18._1,
					_2: _p18._2
				};
			});
		var newTaggers = A3(_elm_lang$core$List$foldl, _elm_lang$core$Time$addMySub, _elm_lang$core$Dict$empty, subs);
		var _p19 = A6(
			_elm_lang$core$Dict$merge,
			leftStep,
			bothStep,
			rightStep,
			newTaggers,
			_p10.processes,
			{
				ctor: '_Tuple3',
				_0: _elm_lang$core$Native_List.fromArray(
					[]),
				_1: _elm_lang$core$Dict$empty,
				_2: _elm_lang$core$Task$succeed(
					{ctor: '_Tuple0'})
			});
		var spawnList = _p19._0;
		var existingDict = _p19._1;
		var killTask = _p19._2;
		return A2(
			_elm_lang$core$Task$andThen,
			killTask,
			function (_p20) {
				return A2(
					_elm_lang$core$Task$andThen,
					A3(_elm_lang$core$Time$spawnHelp, router, spawnList, existingDict),
					function (newProcesses) {
						return _elm_lang$core$Task$succeed(
							A2(_elm_lang$core$Time$State, newTaggers, newProcesses));
					});
			});
	});
var _elm_lang$core$Time$Every = F2(
	function (a, b) {
		return {ctor: 'Every', _0: a, _1: b};
	});
var _elm_lang$core$Time$every = F2(
	function (interval, tagger) {
		return _elm_lang$core$Time$subscription(
			A2(_elm_lang$core$Time$Every, interval, tagger));
	});
var _elm_lang$core$Time$subMap = F2(
	function (f, _p21) {
		var _p22 = _p21;
		return A2(
			_elm_lang$core$Time$Every,
			_p22._0,
			function (_p23) {
				return f(
					_p22._1(_p23));
			});
	});
_elm_lang$core$Native_Platform.effectManagers['Time'] = {pkg: 'elm-lang/core', init: _elm_lang$core$Time$init, onEffects: _elm_lang$core$Time$onEffects, onSelfMsg: _elm_lang$core$Time$onSelfMsg, tag: 'sub', subMap: _elm_lang$core$Time$subMap};

//import Maybe, Native.Array, Native.List, Native.Utils, Result //

var _elm_lang$core$Native_Json = function() {


// CORE DECODERS

function succeed(msg)
{
	return {
		ctor: '<decoder>',
		tag: 'succeed',
		msg: msg
	};
}

function fail(msg)
{
	return {
		ctor: '<decoder>',
		tag: 'fail',
		msg: msg
	};
}

function decodePrimitive(tag)
{
	return {
		ctor: '<decoder>',
		tag: tag
	};
}

function decodeContainer(tag, decoder)
{
	return {
		ctor: '<decoder>',
		tag: tag,
		decoder: decoder
	};
}

function decodeNull(value)
{
	return {
		ctor: '<decoder>',
		tag: 'null',
		value: value
	};
}

function decodeField(field, decoder)
{
	return {
		ctor: '<decoder>',
		tag: 'field',
		field: field,
		decoder: decoder
	};
}

function decodeKeyValuePairs(decoder)
{
	return {
		ctor: '<decoder>',
		tag: 'key-value',
		decoder: decoder
	};
}

function decodeObject(f, decoders)
{
	return {
		ctor: '<decoder>',
		tag: 'map-many',
		func: f,
		decoders: decoders
	};
}

function decodeTuple(f, decoders)
{
	return {
		ctor: '<decoder>',
		tag: 'tuple',
		func: f,
		decoders: decoders
	};
}

function andThen(decoder, callback)
{
	return {
		ctor: '<decoder>',
		tag: 'andThen',
		decoder: decoder,
		callback: callback
	};
}

function customAndThen(decoder, callback)
{
	return {
		ctor: '<decoder>',
		tag: 'customAndThen',
		decoder: decoder,
		callback: callback
	};
}

function oneOf(decoders)
{
	return {
		ctor: '<decoder>',
		tag: 'oneOf',
		decoders: decoders
	};
}


// DECODING OBJECTS

function decodeObject1(f, d1)
{
	return decodeObject(f, [d1]);
}

function decodeObject2(f, d1, d2)
{
	return decodeObject(f, [d1, d2]);
}

function decodeObject3(f, d1, d2, d3)
{
	return decodeObject(f, [d1, d2, d3]);
}

function decodeObject4(f, d1, d2, d3, d4)
{
	return decodeObject(f, [d1, d2, d3, d4]);
}

function decodeObject5(f, d1, d2, d3, d4, d5)
{
	return decodeObject(f, [d1, d2, d3, d4, d5]);
}

function decodeObject6(f, d1, d2, d3, d4, d5, d6)
{
	return decodeObject(f, [d1, d2, d3, d4, d5, d6]);
}

function decodeObject7(f, d1, d2, d3, d4, d5, d6, d7)
{
	return decodeObject(f, [d1, d2, d3, d4, d5, d6, d7]);
}

function decodeObject8(f, d1, d2, d3, d4, d5, d6, d7, d8)
{
	return decodeObject(f, [d1, d2, d3, d4, d5, d6, d7, d8]);
}


// DECODING TUPLES

function decodeTuple1(f, d1)
{
	return decodeTuple(f, [d1]);
}

function decodeTuple2(f, d1, d2)
{
	return decodeTuple(f, [d1, d2]);
}

function decodeTuple3(f, d1, d2, d3)
{
	return decodeTuple(f, [d1, d2, d3]);
}

function decodeTuple4(f, d1, d2, d3, d4)
{
	return decodeTuple(f, [d1, d2, d3, d4]);
}

function decodeTuple5(f, d1, d2, d3, d4, d5)
{
	return decodeTuple(f, [d1, d2, d3, d4, d5]);
}

function decodeTuple6(f, d1, d2, d3, d4, d5, d6)
{
	return decodeTuple(f, [d1, d2, d3, d4, d5, d6]);
}

function decodeTuple7(f, d1, d2, d3, d4, d5, d6, d7)
{
	return decodeTuple(f, [d1, d2, d3, d4, d5, d6, d7]);
}

function decodeTuple8(f, d1, d2, d3, d4, d5, d6, d7, d8)
{
	return decodeTuple(f, [d1, d2, d3, d4, d5, d6, d7, d8]);
}


// DECODE HELPERS

function ok(value)
{
	return { tag: 'ok', value: value };
}

function badPrimitive(type, value)
{
	return { tag: 'primitive', type: type, value: value };
}

function badIndex(index, nestedProblems)
{
	return { tag: 'index', index: index, rest: nestedProblems };
}

function badField(field, nestedProblems)
{
	return { tag: 'field', field: field, rest: nestedProblems };
}

function badOneOf(problems)
{
	return { tag: 'oneOf', problems: problems };
}

function badCustom(msg)
{
	return { tag: 'custom', msg: msg };
}

function bad(msg)
{
	return { tag: 'fail', msg: msg };
}

function badToString(problem)
{
	var context = '_';
	while (problem)
	{
		switch (problem.tag)
		{
			case 'primitive':
				return 'Expecting ' + problem.type
					+ (context === '_' ? '' : ' at ' + context)
					+ ' but instead got: ' + jsToString(problem.value);

			case 'index':
				context += '[' + problem.index + ']';
				problem = problem.rest;
				break;

			case 'field':
				context += '.' + problem.field;
				problem = problem.rest;
				break;

			case 'oneOf':
				var problems = problem.problems;
				for (var i = 0; i < problems.length; i++)
				{
					problems[i] = badToString(problems[i]);
				}
				return 'I ran into the following problems'
					+ (context === '_' ? '' : ' at ' + context)
					+ ':\n\n' + problems.join('\n');

			case 'custom':
				return 'A `customDecoder` failed'
					+ (context === '_' ? '' : ' at ' + context)
					+ ' with the message: ' + problem.msg;

			case 'fail':
				return 'I ran into a `fail` decoder'
					+ (context === '_' ? '' : ' at ' + context)
					+ ': ' + problem.msg;
		}
	}
}

function jsToString(value)
{
	return value === undefined
		? 'undefined'
		: JSON.stringify(value);
}


// DECODE

function runOnString(decoder, string)
{
	var json;
	try
	{
		json = JSON.parse(string);
	}
	catch (e)
	{
		return _elm_lang$core$Result$Err('Given an invalid JSON: ' + e.message);
	}
	return run(decoder, json);
}

function run(decoder, value)
{
	var result = runHelp(decoder, value);
	return (result.tag === 'ok')
		? _elm_lang$core$Result$Ok(result.value)
		: _elm_lang$core$Result$Err(badToString(result));
}

function runHelp(decoder, value)
{
	switch (decoder.tag)
	{
		case 'bool':
			return (typeof value === 'boolean')
				? ok(value)
				: badPrimitive('a Bool', value);

		case 'int':
			if (typeof value !== 'number') {
				return badPrimitive('an Int', value);
			}

			if (-2147483647 < value && value < 2147483647 && (value | 0) === value) {
				return ok(value);
			}

			if (isFinite(value) && !(value % 1)) {
				return ok(value);
			}

			return badPrimitive('an Int', value);

		case 'float':
			return (typeof value === 'number')
				? ok(value)
				: badPrimitive('a Float', value);

		case 'string':
			return (typeof value === 'string')
				? ok(value)
				: (value instanceof String)
					? ok(value + '')
					: badPrimitive('a String', value);

		case 'null':
			return (value === null)
				? ok(decoder.value)
				: badPrimitive('null', value);

		case 'value':
			return ok(value);

		case 'list':
			if (!(value instanceof Array))
			{
				return badPrimitive('a List', value);
			}

			var list = _elm_lang$core$Native_List.Nil;
			for (var i = value.length; i--; )
			{
				var result = runHelp(decoder.decoder, value[i]);
				if (result.tag !== 'ok')
				{
					return badIndex(i, result)
				}
				list = _elm_lang$core$Native_List.Cons(result.value, list);
			}
			return ok(list);

		case 'array':
			if (!(value instanceof Array))
			{
				return badPrimitive('an Array', value);
			}

			var len = value.length;
			var array = new Array(len);
			for (var i = len; i--; )
			{
				var result = runHelp(decoder.decoder, value[i]);
				if (result.tag !== 'ok')
				{
					return badIndex(i, result);
				}
				array[i] = result.value;
			}
			return ok(_elm_lang$core$Native_Array.fromJSArray(array));

		case 'maybe':
			var result = runHelp(decoder.decoder, value);
			return (result.tag === 'ok')
				? ok(_elm_lang$core$Maybe$Just(result.value))
				: ok(_elm_lang$core$Maybe$Nothing);

		case 'field':
			var field = decoder.field;
			if (typeof value !== 'object' || value === null || !(field in value))
			{
				return badPrimitive('an object with a field named `' + field + '`', value);
			}

			var result = runHelp(decoder.decoder, value[field]);
			return (result.tag === 'ok')
				? result
				: badField(field, result);

		case 'key-value':
			if (typeof value !== 'object' || value === null || value instanceof Array)
			{
				return badPrimitive('an object', value);
			}

			var keyValuePairs = _elm_lang$core$Native_List.Nil;
			for (var key in value)
			{
				var result = runHelp(decoder.decoder, value[key]);
				if (result.tag !== 'ok')
				{
					return badField(key, result);
				}
				var pair = _elm_lang$core$Native_Utils.Tuple2(key, result.value);
				keyValuePairs = _elm_lang$core$Native_List.Cons(pair, keyValuePairs);
			}
			return ok(keyValuePairs);

		case 'map-many':
			var answer = decoder.func;
			var decoders = decoder.decoders;
			for (var i = 0; i < decoders.length; i++)
			{
				var result = runHelp(decoders[i], value);
				if (result.tag !== 'ok')
				{
					return result;
				}
				answer = answer(result.value);
			}
			return ok(answer);

		case 'tuple':
			var decoders = decoder.decoders;
			var len = decoders.length;

			if ( !(value instanceof Array) || value.length !== len )
			{
				return badPrimitive('a Tuple with ' + len + ' entries', value);
			}

			var answer = decoder.func;
			for (var i = 0; i < len; i++)
			{
				var result = runHelp(decoders[i], value[i]);
				if (result.tag !== 'ok')
				{
					return badIndex(i, result);
				}
				answer = answer(result.value);
			}
			return ok(answer);

		case 'customAndThen':
			var result = runHelp(decoder.decoder, value);
			if (result.tag !== 'ok')
			{
				return result;
			}
			var realResult = decoder.callback(result.value);
			if (realResult.ctor === 'Err')
			{
				return badCustom(realResult._0);
			}
			return ok(realResult._0);

		case 'andThen':
			var result = runHelp(decoder.decoder, value);
			return (result.tag !== 'ok')
				? result
				: runHelp(decoder.callback(result.value), value);

		case 'oneOf':
			var errors = [];
			var temp = decoder.decoders;
			while (temp.ctor !== '[]')
			{
				var result = runHelp(temp._0, value);

				if (result.tag === 'ok')
				{
					return result;
				}

				errors.push(result);

				temp = temp._1;
			}
			return badOneOf(errors);

		case 'fail':
			return bad(decoder.msg);

		case 'succeed':
			return ok(decoder.msg);
	}
}


// EQUALITY

function equality(a, b)
{
	if (a === b)
	{
		return true;
	}

	if (a.tag !== b.tag)
	{
		return false;
	}

	switch (a.tag)
	{
		case 'succeed':
		case 'fail':
			return a.msg === b.msg;

		case 'bool':
		case 'int':
		case 'float':
		case 'string':
		case 'value':
			return true;

		case 'null':
			return a.value === b.value;

		case 'list':
		case 'array':
		case 'maybe':
		case 'key-value':
			return equality(a.decoder, b.decoder);

		case 'field':
			return a.field === b.field && equality(a.decoder, b.decoder);

		case 'map-many':
		case 'tuple':
			if (a.func !== b.func)
			{
				return false;
			}
			return listEquality(a.decoders, b.decoders);

		case 'andThen':
		case 'customAndThen':
			return a.callback === b.callback && equality(a.decoder, b.decoder);

		case 'oneOf':
			return listEquality(a.decoders, b.decoders);
	}
}

function listEquality(aDecoders, bDecoders)
{
	var len = aDecoders.length;
	if (len !== bDecoders.length)
	{
		return false;
	}
	for (var i = 0; i < len; i++)
	{
		if (!equality(aDecoders[i], bDecoders[i]))
		{
			return false;
		}
	}
	return true;
}


// ENCODE

function encode(indentLevel, value)
{
	return JSON.stringify(value, null, indentLevel);
}

function identity(value)
{
	return value;
}

function encodeObject(keyValuePairs)
{
	var obj = {};
	while (keyValuePairs.ctor !== '[]')
	{
		var pair = keyValuePairs._0;
		obj[pair._0] = pair._1;
		keyValuePairs = keyValuePairs._1;
	}
	return obj;
}

return {
	encode: F2(encode),
	runOnString: F2(runOnString),
	run: F2(run),

	decodeNull: decodeNull,
	decodePrimitive: decodePrimitive,
	decodeContainer: F2(decodeContainer),

	decodeField: F2(decodeField),

	decodeObject1: F2(decodeObject1),
	decodeObject2: F3(decodeObject2),
	decodeObject3: F4(decodeObject3),
	decodeObject4: F5(decodeObject4),
	decodeObject5: F6(decodeObject5),
	decodeObject6: F7(decodeObject6),
	decodeObject7: F8(decodeObject7),
	decodeObject8: F9(decodeObject8),
	decodeKeyValuePairs: decodeKeyValuePairs,

	decodeTuple1: F2(decodeTuple1),
	decodeTuple2: F3(decodeTuple2),
	decodeTuple3: F4(decodeTuple3),
	decodeTuple4: F5(decodeTuple4),
	decodeTuple5: F6(decodeTuple5),
	decodeTuple6: F7(decodeTuple6),
	decodeTuple7: F8(decodeTuple7),
	decodeTuple8: F9(decodeTuple8),

	andThen: F2(andThen),
	customAndThen: F2(customAndThen),
	fail: fail,
	succeed: succeed,
	oneOf: oneOf,

	identity: identity,
	encodeNull: null,
	encodeArray: _elm_lang$core$Native_Array.toJSArray,
	encodeList: _elm_lang$core$Native_List.toArray,
	encodeObject: encodeObject,

	equality: equality
};

}();

var _elm_lang$core$Json_Encode$list = _elm_lang$core$Native_Json.encodeList;
var _elm_lang$core$Json_Encode$array = _elm_lang$core$Native_Json.encodeArray;
var _elm_lang$core$Json_Encode$object = _elm_lang$core$Native_Json.encodeObject;
var _elm_lang$core$Json_Encode$null = _elm_lang$core$Native_Json.encodeNull;
var _elm_lang$core$Json_Encode$bool = _elm_lang$core$Native_Json.identity;
var _elm_lang$core$Json_Encode$float = _elm_lang$core$Native_Json.identity;
var _elm_lang$core$Json_Encode$int = _elm_lang$core$Native_Json.identity;
var _elm_lang$core$Json_Encode$string = _elm_lang$core$Native_Json.identity;
var _elm_lang$core$Json_Encode$encode = _elm_lang$core$Native_Json.encode;
var _elm_lang$core$Json_Encode$Value = {ctor: 'Value'};

var _elm_lang$core$Json_Decode$tuple8 = _elm_lang$core$Native_Json.decodeTuple8;
var _elm_lang$core$Json_Decode$tuple7 = _elm_lang$core$Native_Json.decodeTuple7;
var _elm_lang$core$Json_Decode$tuple6 = _elm_lang$core$Native_Json.decodeTuple6;
var _elm_lang$core$Json_Decode$tuple5 = _elm_lang$core$Native_Json.decodeTuple5;
var _elm_lang$core$Json_Decode$tuple4 = _elm_lang$core$Native_Json.decodeTuple4;
var _elm_lang$core$Json_Decode$tuple3 = _elm_lang$core$Native_Json.decodeTuple3;
var _elm_lang$core$Json_Decode$tuple2 = _elm_lang$core$Native_Json.decodeTuple2;
var _elm_lang$core$Json_Decode$tuple1 = _elm_lang$core$Native_Json.decodeTuple1;
var _elm_lang$core$Json_Decode$succeed = _elm_lang$core$Native_Json.succeed;
var _elm_lang$core$Json_Decode$fail = _elm_lang$core$Native_Json.fail;
var _elm_lang$core$Json_Decode$andThen = _elm_lang$core$Native_Json.andThen;
var _elm_lang$core$Json_Decode$customDecoder = _elm_lang$core$Native_Json.customAndThen;
var _elm_lang$core$Json_Decode$decodeValue = _elm_lang$core$Native_Json.run;
var _elm_lang$core$Json_Decode$value = _elm_lang$core$Native_Json.decodePrimitive('value');
var _elm_lang$core$Json_Decode$maybe = function (decoder) {
	return A2(_elm_lang$core$Native_Json.decodeContainer, 'maybe', decoder);
};
var _elm_lang$core$Json_Decode$null = _elm_lang$core$Native_Json.decodeNull;
var _elm_lang$core$Json_Decode$array = function (decoder) {
	return A2(_elm_lang$core$Native_Json.decodeContainer, 'array', decoder);
};
var _elm_lang$core$Json_Decode$list = function (decoder) {
	return A2(_elm_lang$core$Native_Json.decodeContainer, 'list', decoder);
};
var _elm_lang$core$Json_Decode$bool = _elm_lang$core$Native_Json.decodePrimitive('bool');
var _elm_lang$core$Json_Decode$int = _elm_lang$core$Native_Json.decodePrimitive('int');
var _elm_lang$core$Json_Decode$float = _elm_lang$core$Native_Json.decodePrimitive('float');
var _elm_lang$core$Json_Decode$string = _elm_lang$core$Native_Json.decodePrimitive('string');
var _elm_lang$core$Json_Decode$oneOf = _elm_lang$core$Native_Json.oneOf;
var _elm_lang$core$Json_Decode$keyValuePairs = _elm_lang$core$Native_Json.decodeKeyValuePairs;
var _elm_lang$core$Json_Decode$object8 = _elm_lang$core$Native_Json.decodeObject8;
var _elm_lang$core$Json_Decode$object7 = _elm_lang$core$Native_Json.decodeObject7;
var _elm_lang$core$Json_Decode$object6 = _elm_lang$core$Native_Json.decodeObject6;
var _elm_lang$core$Json_Decode$object5 = _elm_lang$core$Native_Json.decodeObject5;
var _elm_lang$core$Json_Decode$object4 = _elm_lang$core$Native_Json.decodeObject4;
var _elm_lang$core$Json_Decode$object3 = _elm_lang$core$Native_Json.decodeObject3;
var _elm_lang$core$Json_Decode$object2 = _elm_lang$core$Native_Json.decodeObject2;
var _elm_lang$core$Json_Decode$object1 = _elm_lang$core$Native_Json.decodeObject1;
var _elm_lang$core$Json_Decode_ops = _elm_lang$core$Json_Decode_ops || {};
_elm_lang$core$Json_Decode_ops[':='] = _elm_lang$core$Native_Json.decodeField;
var _elm_lang$core$Json_Decode$at = F2(
	function (fields, decoder) {
		return A3(
			_elm_lang$core$List$foldr,
			F2(
				function (x, y) {
					return A2(_elm_lang$core$Json_Decode_ops[':='], x, y);
				}),
			decoder,
			fields);
	});
var _elm_lang$core$Json_Decode$decodeString = _elm_lang$core$Native_Json.runOnString;
var _elm_lang$core$Json_Decode$map = _elm_lang$core$Native_Json.decodeObject1;
var _elm_lang$core$Json_Decode$dict = function (decoder) {
	return A2(
		_elm_lang$core$Json_Decode$map,
		_elm_lang$core$Dict$fromList,
		_elm_lang$core$Json_Decode$keyValuePairs(decoder));
};
var _elm_lang$core$Json_Decode$Decoder = {ctor: 'Decoder'};

//import Native.Json //

var _elm_lang$virtual_dom$Native_VirtualDom = function() {

var STYLE_KEY = 'STYLE';
var EVENT_KEY = 'EVENT';
var ATTR_KEY = 'ATTR';
var ATTR_NS_KEY = 'ATTR_NS';



////////////  VIRTUAL DOM NODES  ////////////


function text(string)
{
	return {
		type: 'text',
		text: string
	};
}


function node(tag)
{
	return F2(function(factList, kidList) {
		return nodeHelp(tag, factList, kidList);
	});
}


function nodeHelp(tag, factList, kidList)
{
	var organized = organizeFacts(factList);
	var namespace = organized.namespace;
	var facts = organized.facts;

	var children = [];
	var descendantsCount = 0;
	while (kidList.ctor !== '[]')
	{
		var kid = kidList._0;
		descendantsCount += (kid.descendantsCount || 0);
		children.push(kid);
		kidList = kidList._1;
	}
	descendantsCount += children.length;

	return {
		type: 'node',
		tag: tag,
		facts: facts,
		children: children,
		namespace: namespace,
		descendantsCount: descendantsCount
	};
}


function keyedNode(tag, factList, kidList)
{
	var organized = organizeFacts(factList);
	var namespace = organized.namespace;
	var facts = organized.facts;

	var children = [];
	var descendantsCount = 0;
	while (kidList.ctor !== '[]')
	{
		var kid = kidList._0;
		descendantsCount += (kid._1.descendantsCount || 0);
		children.push(kid);
		kidList = kidList._1;
	}
	descendantsCount += children.length;

	return {
		type: 'keyed-node',
		tag: tag,
		facts: facts,
		children: children,
		namespace: namespace,
		descendantsCount: descendantsCount
	};
}


function custom(factList, model, impl)
{
	var facts = organizeFacts(factList).facts;

	return {
		type: 'custom',
		facts: facts,
		model: model,
		impl: impl
	};
}


function map(tagger, node)
{
	return {
		type: 'tagger',
		tagger: tagger,
		node: node,
		descendantsCount: 1 + (node.descendantsCount || 0)
	};
}


function thunk(func, args, thunk)
{
	return {
		type: 'thunk',
		func: func,
		args: args,
		thunk: thunk,
		node: undefined
	};
}

function lazy(fn, a)
{
	return thunk(fn, [a], function() {
		return fn(a);
	});
}

function lazy2(fn, a, b)
{
	return thunk(fn, [a,b], function() {
		return A2(fn, a, b);
	});
}

function lazy3(fn, a, b, c)
{
	return thunk(fn, [a,b,c], function() {
		return A3(fn, a, b, c);
	});
}



// FACTS


function organizeFacts(factList)
{
	var namespace, facts = {};

	while (factList.ctor !== '[]')
	{
		var entry = factList._0;
		var key = entry.key;

		if (key === ATTR_KEY || key === ATTR_NS_KEY || key === EVENT_KEY)
		{
			var subFacts = facts[key] || {};
			subFacts[entry.realKey] = entry.value;
			facts[key] = subFacts;
		}
		else if (key === STYLE_KEY)
		{
			var styles = facts[key] || {};
			var styleList = entry.value;
			while (styleList.ctor !== '[]')
			{
				var style = styleList._0;
				styles[style._0] = style._1;
				styleList = styleList._1;
			}
			facts[key] = styles;
		}
		else if (key === 'namespace')
		{
			namespace = entry.value;
		}
		else
		{
			facts[key] = entry.value;
		}
		factList = factList._1;
	}

	return {
		facts: facts,
		namespace: namespace
	};
}



////////////  PROPERTIES AND ATTRIBUTES  ////////////


function style(value)
{
	return {
		key: STYLE_KEY,
		value: value
	};
}


function property(key, value)
{
	return {
		key: key,
		value: value
	};
}


function attribute(key, value)
{
	return {
		key: ATTR_KEY,
		realKey: key,
		value: value
	};
}


function attributeNS(namespace, key, value)
{
	return {
		key: ATTR_NS_KEY,
		realKey: key,
		value: {
			value: value,
			namespace: namespace
		}
	};
}


function on(name, options, decoder)
{
	return {
		key: EVENT_KEY,
		realKey: name,
		value: {
			options: options,
			decoder: decoder
		}
	};
}


function equalEvents(a, b)
{
	if (!a.options === b.options)
	{
		if (a.stopPropagation !== b.stopPropagation || a.preventDefault !== b.preventDefault)
		{
			return false;
		}
	}
	return _elm_lang$core$Native_Json.equality(a.decoder, b.decoder);
}



////////////  RENDERER  ////////////


function renderer(parent, tagger, initialVirtualNode)
{
	var eventNode = { tagger: tagger, parent: undefined };

	var domNode = render(initialVirtualNode, eventNode);
	parent.appendChild(domNode);

	var state = 'NO_REQUEST';
	var currentVirtualNode = initialVirtualNode;
	var nextVirtualNode = initialVirtualNode;

	function registerVirtualNode(vNode)
	{
		if (state === 'NO_REQUEST')
		{
			rAF(updateIfNeeded);
		}
		state = 'PENDING_REQUEST';
		nextVirtualNode = vNode;
	}

	function updateIfNeeded()
	{
		switch (state)
		{
			case 'NO_REQUEST':
				throw new Error(
					'Unexpected draw callback.\n' +
					'Please report this to <https://github.com/elm-lang/core/issues>.'
				);

			case 'PENDING_REQUEST':
				rAF(updateIfNeeded);
				state = 'EXTRA_REQUEST';

				var patches = diff(currentVirtualNode, nextVirtualNode);
				domNode = applyPatches(domNode, currentVirtualNode, patches, eventNode);
				currentVirtualNode = nextVirtualNode;

				return;

			case 'EXTRA_REQUEST':
				state = 'NO_REQUEST';
				return;
		}
	}

	return { update: registerVirtualNode };
}


var rAF =
	typeof requestAnimationFrame !== 'undefined'
		? requestAnimationFrame
		: function(cb) { setTimeout(cb, 1000 / 60); };



////////////  RENDER  ////////////


function render(vNode, eventNode)
{
	switch (vNode.type)
	{
		case 'thunk':
			if (!vNode.node)
			{
				vNode.node = vNode.thunk();
			}
			return render(vNode.node, eventNode);

		case 'tagger':
			var subNode = vNode.node;
			var tagger = vNode.tagger;

			while (subNode.type === 'tagger')
			{
				typeof tagger !== 'object'
					? tagger = [tagger, subNode.tagger]
					: tagger.push(subNode.tagger);

				subNode = subNode.node;
			}

			var subEventRoot = {
				tagger: tagger,
				parent: eventNode
			};

			var domNode = render(subNode, subEventRoot);
			domNode.elm_event_node_ref = subEventRoot;
			return domNode;

		case 'text':
			return document.createTextNode(vNode.text);

		case 'node':
			var domNode = vNode.namespace
				? document.createElementNS(vNode.namespace, vNode.tag)
				: document.createElement(vNode.tag);

			applyFacts(domNode, eventNode, vNode.facts);

			var children = vNode.children;

			for (var i = 0; i < children.length; i++)
			{
				domNode.appendChild(render(children[i], eventNode));
			}

			return domNode;

		case 'keyed-node':
			var domNode = vNode.namespace
				? document.createElementNS(vNode.namespace, vNode.tag)
				: document.createElement(vNode.tag);

			applyFacts(domNode, eventNode, vNode.facts);

			var children = vNode.children;

			for (var i = 0; i < children.length; i++)
			{
				domNode.appendChild(render(children[i]._1, eventNode));
			}

			return domNode;

		case 'custom':
			var domNode = vNode.impl.render(vNode.model);
			applyFacts(domNode, eventNode, vNode.facts);
			return domNode;
	}
}



////////////  APPLY FACTS  ////////////


function applyFacts(domNode, eventNode, facts)
{
	for (var key in facts)
	{
		var value = facts[key];

		switch (key)
		{
			case STYLE_KEY:
				applyStyles(domNode, value);
				break;

			case EVENT_KEY:
				applyEvents(domNode, eventNode, value);
				break;

			case ATTR_KEY:
				applyAttrs(domNode, value);
				break;

			case ATTR_NS_KEY:
				applyAttrsNS(domNode, value);
				break;

			case 'value':
				if (domNode[key] !== value)
				{
					domNode[key] = value;
				}
				break;

			default:
				domNode[key] = value;
				break;
		}
	}
}

function applyStyles(domNode, styles)
{
	var domNodeStyle = domNode.style;

	for (var key in styles)
	{
		domNodeStyle[key] = styles[key];
	}
}

function applyEvents(domNode, eventNode, events)
{
	var allHandlers = domNode.elm_handlers || {};

	for (var key in events)
	{
		var handler = allHandlers[key];
		var value = events[key];

		if (typeof value === 'undefined')
		{
			domNode.removeEventListener(key, handler);
			allHandlers[key] = undefined;
		}
		else if (typeof handler === 'undefined')
		{
			var handler = makeEventHandler(eventNode, value);
			domNode.addEventListener(key, handler);
			allHandlers[key] = handler;
		}
		else
		{
			handler.info = value;
		}
	}

	domNode.elm_handlers = allHandlers;
}

function makeEventHandler(eventNode, info)
{
	function eventHandler(event)
	{
		var info = eventHandler.info;

		var value = A2(_elm_lang$core$Native_Json.run, info.decoder, event);

		if (value.ctor === 'Ok')
		{
			var options = info.options;
			if (options.stopPropagation)
			{
				event.stopPropagation();
			}
			if (options.preventDefault)
			{
				event.preventDefault();
			}

			var message = value._0;

			var currentEventNode = eventNode;
			while (currentEventNode)
			{
				var tagger = currentEventNode.tagger;
				if (typeof tagger === 'function')
				{
					message = tagger(message);
				}
				else
				{
					for (var i = tagger.length; i--; )
					{
						message = tagger[i](message);
					}
				}
				currentEventNode = currentEventNode.parent;
			}
		}
	};

	eventHandler.info = info;

	return eventHandler;
}

function applyAttrs(domNode, attrs)
{
	for (var key in attrs)
	{
		var value = attrs[key];
		if (typeof value === 'undefined')
		{
			domNode.removeAttribute(key);
		}
		else
		{
			domNode.setAttribute(key, value);
		}
	}
}

function applyAttrsNS(domNode, nsAttrs)
{
	for (var key in nsAttrs)
	{
		var pair = nsAttrs[key];
		var namespace = pair.namespace;
		var value = pair.value;

		if (typeof value === 'undefined')
		{
			domNode.removeAttributeNS(namespace, key);
		}
		else
		{
			domNode.setAttributeNS(namespace, key, value);
		}
	}
}



////////////  DIFF  ////////////


function diff(a, b)
{
	var patches = [];
	diffHelp(a, b, patches, 0);
	return patches;
}


function makePatch(type, index, data)
{
	return {
		index: index,
		type: type,
		data: data,
		domNode: undefined,
		eventNode: undefined
	};
}


function diffHelp(a, b, patches, index)
{
	if (a === b)
	{
		return;
	}

	var aType = a.type;
	var bType = b.type;

	// Bail if you run into different types of nodes. Implies that the
	// structure has changed significantly and it's not worth a diff.
	if (aType !== bType)
	{
		patches.push(makePatch('p-redraw', index, b));
		return;
	}

	// Now we know that both nodes are the same type.
	switch (bType)
	{
		case 'thunk':
			var aArgs = a.args;
			var bArgs = b.args;
			var i = aArgs.length;
			var same = a.func === b.func && i === bArgs.length;
			while (same && i--)
			{
				same = aArgs[i] === bArgs[i];
			}
			if (same)
			{
				b.node = a.node;
				return;
			}
			b.node = b.thunk();
			var subPatches = [];
			diffHelp(a.node, b.node, subPatches, 0);
			if (subPatches.length > 0)
			{
				patches.push(makePatch('p-thunk', index, subPatches));
			}
			return;

		case 'tagger':
			// gather nested taggers
			var aTaggers = a.tagger;
			var bTaggers = b.tagger;
			var nesting = false;

			var aSubNode = a.node;
			while (aSubNode.type === 'tagger')
			{
				nesting = true;

				typeof aTaggers !== 'object'
					? aTaggers = [aTaggers, aSubNode.tagger]
					: aTaggers.push(aSubNode.tagger);

				aSubNode = aSubNode.node;
			}

			var bSubNode = b.node;
			while (bSubNode.type === 'tagger')
			{
				nesting = true;

				typeof bTaggers !== 'object'
					? bTaggers = [bTaggers, bSubNode.tagger]
					: bTaggers.push(bSubNode.tagger);

				bSubNode = bSubNode.node;
			}

			// Just bail if different numbers of taggers. This implies the
			// structure of the virtual DOM has changed.
			if (nesting && aTaggers.length !== bTaggers.length)
			{
				patches.push(makePatch('p-redraw', index, b));
				return;
			}

			// check if taggers are "the same"
			if (nesting ? !pairwiseRefEqual(aTaggers, bTaggers) : aTaggers !== bTaggers)
			{
				patches.push(makePatch('p-tagger', index, bTaggers));
			}

			// diff everything below the taggers
			diffHelp(aSubNode, bSubNode, patches, index + 1);
			return;

		case 'text':
			if (a.text !== b.text)
			{
				patches.push(makePatch('p-text', index, b.text));
				return;
			}

			return;

		case 'node':
			// Bail if obvious indicators have changed. Implies more serious
			// structural changes such that it's not worth it to diff.
			if (a.tag !== b.tag || a.namespace !== b.namespace)
			{
				patches.push(makePatch('p-redraw', index, b));
				return;
			}

			var factsDiff = diffFacts(a.facts, b.facts);

			if (typeof factsDiff !== 'undefined')
			{
				patches.push(makePatch('p-facts', index, factsDiff));
			}

			diffChildren(a, b, patches, index);
			return;

		case 'keyed-node':
			// Bail if obvious indicators have changed. Implies more serious
			// structural changes such that it's not worth it to diff.
			if (a.tag !== b.tag || a.namespace !== b.namespace)
			{
				patches.push(makePatch('p-redraw', index, b));
				return;
			}

			var factsDiff = diffFacts(a.facts, b.facts);

			if (typeof factsDiff !== 'undefined')
			{
				patches.push(makePatch('p-facts', index, factsDiff));
			}

			diffKeyedChildren(a, b, patches, index);
			return;

		case 'custom':
			if (a.impl !== b.impl)
			{
				patches.push(makePatch('p-redraw', index, b));
				return;
			}

			var factsDiff = diffFacts(a.facts, b.facts);
			if (typeof factsDiff !== 'undefined')
			{
				patches.push(makePatch('p-facts', index, factsDiff));
			}

			var patch = b.impl.diff(a,b);
			if (patch)
			{
				patches.push(makePatch('p-custom', index, patch));
				return;
			}

			return;
	}
}


// assumes the incoming arrays are the same length
function pairwiseRefEqual(as, bs)
{
	for (var i = 0; i < as.length; i++)
	{
		if (as[i] !== bs[i])
		{
			return false;
		}
	}

	return true;
}


// TODO Instead of creating a new diff object, it's possible to just test if
// there *is* a diff. During the actual patch, do the diff again and make the
// modifications directly. This way, there's no new allocations. Worth it?
function diffFacts(a, b, category)
{
	var diff;

	// look for changes and removals
	for (var aKey in a)
	{
		if (aKey === STYLE_KEY || aKey === EVENT_KEY || aKey === ATTR_KEY || aKey === ATTR_NS_KEY)
		{
			var subDiff = diffFacts(a[aKey], b[aKey] || {}, aKey);
			if (subDiff)
			{
				diff = diff || {};
				diff[aKey] = subDiff;
			}
			continue;
		}

		// remove if not in the new facts
		if (!(aKey in b))
		{
			diff = diff || {};
			diff[aKey] =
				(typeof category === 'undefined')
					? (typeof a[aKey] === 'string' ? '' : null)
					:
				(category === STYLE_KEY)
					? ''
					:
				(category === EVENT_KEY || category === ATTR_KEY)
					? undefined
					:
				{ namespace: a[aKey].namespace, value: undefined };

			continue;
		}

		var aValue = a[aKey];
		var bValue = b[aKey];

		// reference equal, so don't worry about it
		if (aValue === bValue && aKey !== 'value'
			|| category === EVENT_KEY && equalEvents(aValue, bValue))
		{
			continue;
		}

		diff = diff || {};
		diff[aKey] = bValue;
	}

	// add new stuff
	for (var bKey in b)
	{
		if (!(bKey in a))
		{
			diff = diff || {};
			diff[bKey] = b[bKey];
		}
	}

	return diff;
}


function diffChildren(aParent, bParent, patches, rootIndex)
{
	var aChildren = aParent.children;
	var bChildren = bParent.children;

	var aLen = aChildren.length;
	var bLen = bChildren.length;

	// FIGURE OUT IF THERE ARE INSERTS OR REMOVALS

	if (aLen > bLen)
	{
		patches.push(makePatch('p-remove-last', rootIndex, aLen - bLen));
	}
	else if (aLen < bLen)
	{
		patches.push(makePatch('p-append', rootIndex, bChildren.slice(aLen)));
	}

	// PAIRWISE DIFF EVERYTHING ELSE

	var index = rootIndex;
	var minLen = aLen < bLen ? aLen : bLen;
	for (var i = 0; i < minLen; i++)
	{
		index++;
		var aChild = aChildren[i];
		diffHelp(aChild, bChildren[i], patches, index);
		index += aChild.descendantsCount || 0;
	}
}



////////////  KEYED DIFF  ////////////


function diffKeyedChildren(aParent, bParent, patches, rootIndex)
{
	var localPatches = [];

	var changes = {}; // Dict String Entry
	var inserts = []; // Array { index : Int, entry : Entry }
	// type Entry = { tag : String, vnode : VNode, index : Int, data : _ }

	var aChildren = aParent.children;
	var bChildren = bParent.children;
	var aLen = aChildren.length;
	var bLen = bChildren.length;
	var aIndex = 0;
	var bIndex = 0;

	var index = rootIndex;

	while (aIndex < aLen && bIndex < bLen)
	{
		var a = aChildren[aIndex];
		var b = bChildren[bIndex];

		var aKey = a._0;
		var bKey = b._0;
		var aNode = a._1;
		var bNode = b._1;

		// check if keys match

		if (aKey === bKey)
		{
			index++;
			diffHelp(aNode, bNode, localPatches, index);
			index += aNode.descendantsCount || 0;

			aIndex++;
			bIndex++;
			continue;
		}

		// look ahead 1 to detect insertions and removals.

		var aLookAhead = aIndex + 1 < aLen;
		var bLookAhead = bIndex + 1 < bLen;

		if (aLookAhead)
		{
			var aNext = aChildren[aIndex + 1];
			var aNextKey = aNext._0;
			var aNextNode = aNext._1;
			var oldMatch = bKey === aNextKey;
		}

		if (bLookAhead)
		{
			var bNext = bChildren[bIndex + 1];
			var bNextKey = bNext._0;
			var bNextNode = bNext._1;
			var newMatch = aKey === bNextKey;
		}


		// swap a and b
		if (aLookAhead && bLookAhead && newMatch && oldMatch)
		{
			index++;
			diffHelp(aNode, bNextNode, localPatches, index);
			insertNode(changes, localPatches, aKey, bNode, bIndex, inserts);
			index += aNode.descendantsCount || 0;

			index++;
			removeNode(changes, localPatches, aKey, aNextNode, index);
			index += aNextNode.descendantsCount || 0;

			aIndex += 2;
			bIndex += 2;
			continue;
		}

		// insert b
		if (bLookAhead && newMatch)
		{
			index++;
			insertNode(changes, localPatches, bKey, bNode, bIndex, inserts);
			diffHelp(aNode, bNextNode, localPatches, index);
			index += aNode.descendantsCount || 0;

			aIndex += 1;
			bIndex += 2;
			continue;
		}

		// remove a
		if (aLookAhead && oldMatch)
		{
			index++;
			removeNode(changes, localPatches, aKey, aNode, index);
			index += aNode.descendantsCount || 0;

			index++;
			diffHelp(aNextNode, bNode, localPatches, index);
			index += aNextNode.descendantsCount || 0;

			aIndex += 2;
			bIndex += 1;
			continue;
		}

		// remove a, insert b
		if (aLookAhead && bLookAhead && aNextKey === bNextKey)
		{
			index++;
			removeNode(changes, localPatches, aKey, aNode, index);
			insertNode(changes, localPatches, bKey, bNode, bIndex, inserts);
			index += aNode.descendantsCount || 0;

			index++;
			diffHelp(aNextNode, bNextNode, localPatches, index);
			index += aNextNode.descendantsCount || 0;

			aIndex += 2;
			bIndex += 2;
			continue;
		}

		break;
	}

	// eat up any remaining nodes with removeNode and insertNode

	while (aIndex < aLen)
	{
		index++;
		var a = aChildren[aIndex];
		var aNode = a._1;
		removeNode(changes, localPatches, a._0, aNode, index);
		index += aNode.descendantsCount || 0;
		aIndex++;
	}

	var endInserts;
	while (bIndex < bLen)
	{
		endInserts = endInserts || [];
		var b = bChildren[bIndex];
		insertNode(changes, localPatches, b._0, b._1, undefined, endInserts);
		bIndex++;
	}

	if (localPatches.length > 0 || inserts.length > 0 || typeof endInserts !== 'undefined')
	{
		patches.push(makePatch('p-reorder', rootIndex, {
			patches: localPatches,
			inserts: inserts,
			endInserts: endInserts
		}));
	}
}



////////////  CHANGES FROM KEYED DIFF  ////////////


var POSTFIX = '_elmW6BL';


function insertNode(changes, localPatches, key, vnode, bIndex, inserts)
{
	var entry = changes[key];

	// never seen this key before
	if (typeof entry === 'undefined')
	{
		entry = {
			tag: 'insert',
			vnode: vnode,
			index: bIndex,
			data: undefined
		};

		inserts.push({ index: bIndex, entry: entry });
		changes[key] = entry;

		return;
	}

	// this key was removed earlier, a match!
	if (entry.tag === 'remove')
	{
		inserts.push({ index: bIndex, entry: entry });

		entry.tag = 'move';
		var subPatches = [];
		diffHelp(entry.vnode, vnode, subPatches, entry.index);
		entry.index = bIndex;
		entry.data.data = {
			patches: subPatches,
			entry: entry
		};

		return;
	}

	// this key has already been inserted or moved, a duplicate!
	insertNode(changes, localPatches, key + POSTFIX, vnode, bIndex, inserts);
}


function removeNode(changes, localPatches, key, vnode, index)
{
	var entry = changes[key];

	// never seen this key before
	if (typeof entry === 'undefined')
	{
		var patch = makePatch('p-remove', index, undefined);
		localPatches.push(patch);

		changes[key] = {
			tag: 'remove',
			vnode: vnode,
			index: index,
			data: patch
		};

		return;
	}

	// this key was inserted earlier, a match!
	if (entry.tag === 'insert')
	{
		entry.tag = 'move';
		var subPatches = [];
		diffHelp(vnode, entry.vnode, subPatches, index);

		var patch = makePatch('p-remove', index, {
			patches: subPatches,
			entry: entry
		});
		localPatches.push(patch);

		return;
	}

	// this key has already been removed or moved, a duplicate!
	removeNode(changes, localPatches, key + POSTFIX, vnode, index);
}



////////////  ADD DOM NODES  ////////////
//
// Each DOM node has an "index" assigned in order of traversal. It is important
// to minimize our crawl over the actual DOM, so these indexes (along with the
// descendantsCount of virtual nodes) let us skip touching entire subtrees of
// the DOM if we know there are no patches there.


function addDomNodes(domNode, vNode, patches, eventNode)
{
	addDomNodesHelp(domNode, vNode, patches, 0, 0, vNode.descendantsCount, eventNode);
}


// assumes `patches` is non-empty and indexes increase monotonically.
function addDomNodesHelp(domNode, vNode, patches, i, low, high, eventNode)
{
	var patch = patches[i];
	var index = patch.index;

	while (index === low)
	{
		var patchType = patch.type;

		if (patchType === 'p-thunk')
		{
			addDomNodes(domNode, vNode.node, patch.data, eventNode);
		}
		else if (patchType === 'p-reorder')
		{
			patch.domNode = domNode;
			patch.eventNode = eventNode;

			var subPatches = patch.data.patches;
			if (subPatches.length > 0)
			{
				addDomNodesHelp(domNode, vNode, subPatches, 0, low, high, eventNode);
			}
		}
		else if (patchType === 'p-remove')
		{
			patch.domNode = domNode;
			patch.eventNode = eventNode;

			var data = patch.data;
			if (typeof data !== 'undefined')
			{
				data.entry.data = domNode;
				var subPatches = data.patches;
				if (subPatches.length > 0)
				{
					addDomNodesHelp(domNode, vNode, subPatches, 0, low, high, eventNode);
				}
			}
		}
		else
		{
			patch.domNode = domNode;
			patch.eventNode = eventNode;
		}

		i++;

		if (!(patch = patches[i]) || (index = patch.index) > high)
		{
			return i;
		}
	}

	switch (vNode.type)
	{
		case 'tagger':
			var subNode = vNode.node;

			while (subNode.type === "tagger")
			{
				subNode = subNode.node;
			}

			return addDomNodesHelp(domNode, subNode, patches, i, low + 1, high, domNode.elm_event_node_ref);

		case 'node':
			var vChildren = vNode.children;
			var childNodes = domNode.childNodes;
			for (var j = 0; j < vChildren.length; j++)
			{
				low++;
				var vChild = vChildren[j];
				var nextLow = low + (vChild.descendantsCount || 0);
				if (low <= index && index <= nextLow)
				{
					i = addDomNodesHelp(childNodes[j], vChild, patches, i, low, nextLow, eventNode);
					if (!(patch = patches[i]) || (index = patch.index) > high)
					{
						return i;
					}
				}
				low = nextLow;
			}
			return i;

		case 'keyed-node':
			var vChildren = vNode.children;
			var childNodes = domNode.childNodes;
			for (var j = 0; j < vChildren.length; j++)
			{
				low++;
				var vChild = vChildren[j]._1;
				var nextLow = low + (vChild.descendantsCount || 0);
				if (low <= index && index <= nextLow)
				{
					i = addDomNodesHelp(childNodes[j], vChild, patches, i, low, nextLow, eventNode);
					if (!(patch = patches[i]) || (index = patch.index) > high)
					{
						return i;
					}
				}
				low = nextLow;
			}
			return i;

		case 'text':
		case 'thunk':
			throw new Error('should never traverse `text` or `thunk` nodes like this');
	}
}



////////////  APPLY PATCHES  ////////////


function applyPatches(rootDomNode, oldVirtualNode, patches, eventNode)
{
	if (patches.length === 0)
	{
		return rootDomNode;
	}

	addDomNodes(rootDomNode, oldVirtualNode, patches, eventNode);
	return applyPatchesHelp(rootDomNode, patches);
}

function applyPatchesHelp(rootDomNode, patches)
{
	for (var i = 0; i < patches.length; i++)
	{
		var patch = patches[i];
		var localDomNode = patch.domNode
		var newNode = applyPatch(localDomNode, patch);
		if (localDomNode === rootDomNode)
		{
			rootDomNode = newNode;
		}
	}
	return rootDomNode;
}

function applyPatch(domNode, patch)
{
	switch (patch.type)
	{
		case 'p-redraw':
			return redraw(domNode, patch.data, patch.eventNode);

		case 'p-facts':
			applyFacts(domNode, patch.eventNode, patch.data);
			return domNode;

		case 'p-text':
			domNode.replaceData(0, domNode.length, patch.data);
			return domNode;

		case 'p-thunk':
			return applyPatchesHelp(domNode, patch.data);

		case 'p-tagger':
			domNode.elm_event_node_ref.tagger = patch.data;
			return domNode;

		case 'p-remove-last':
			var i = patch.data;
			while (i--)
			{
				domNode.removeChild(domNode.lastChild);
			}
			return domNode;

		case 'p-append':
			var newNodes = patch.data;
			for (var i = 0; i < newNodes.length; i++)
			{
				domNode.appendChild(render(newNodes[i], patch.eventNode));
			}
			return domNode;

		case 'p-remove':
			var data = patch.data;
			if (typeof data === 'undefined')
			{
				domNode.parentNode.removeChild(domNode);
				return domNode;
			}
			var entry = data.entry;
			if (typeof entry.index !== 'undefined')
			{
				domNode.parentNode.removeChild(domNode);
			}
			entry.data = applyPatchesHelp(domNode, data.patches);
			return domNode;

		case 'p-reorder':
			var data = patch.data;

			// end inserts
			var endInserts = data.endInserts;
			var end;
			if (typeof endInserts !== 'undefined')
			{
				if (endInserts.length === 1)
				{
					var insert = endInserts[0];
					var entry = insert.entry;
					var end = entry.tag === 'move'
						? entry.data
						: render(entry.vnode, patch.eventNode);
				}
				else
				{
					end = document.createDocumentFragment();
					for (var i = 0; i < endInserts.length; i++)
					{
						var insert = endInserts[i];
						var entry = insert.entry;
						var node = entry.tag === 'move'
							? entry.data
							: render(entry.vnode, patch.eventNode);
						end.appendChild(node);
					}
				}
			}

			// removals
			domNode = applyPatchesHelp(domNode, data.patches);

			// inserts
			var inserts = data.inserts;
			for (var i = 0; i < inserts.length; i++)
			{
				var insert = inserts[i];
				var entry = insert.entry;
				var node = entry.tag === 'move'
					? entry.data
					: render(entry.vnode, patch.eventNode);
				domNode.insertBefore(node, domNode.childNodes[insert.index]);
			}

			if (typeof end !== 'undefined')
			{
				domNode.appendChild(end);
			}

			return domNode;

		case 'p-custom':
			var impl = patch.data;
			return impl.applyPatch(domNode, impl.data);

		default:
			throw new Error('Ran into an unknown patch!');
	}
}


function redraw(domNode, vNode, eventNode)
{
	var parentNode = domNode.parentNode;
	var newNode = render(vNode, eventNode);

	if (typeof newNode.elm_event_node_ref === 'undefined')
	{
		newNode.elm_event_node_ref = domNode.elm_event_node_ref;
	}

	if (parentNode && newNode !== domNode)
	{
		parentNode.replaceChild(newNode, domNode);
	}
	return newNode;
}



////////////  PROGRAMS  ////////////


function programWithFlags(details)
{
	return {
		init: details.init,
		update: details.update,
		subscriptions: details.subscriptions,
		view: details.view,
		renderer: renderer
	};
}


return {
	node: node,
	text: text,

	custom: custom,

	map: F2(map),

	on: F3(on),
	style: style,
	property: F2(property),
	attribute: F2(attribute),
	attributeNS: F3(attributeNS),

	lazy: F2(lazy),
	lazy2: F3(lazy2),
	lazy3: F4(lazy3),
	keyedNode: F3(keyedNode),

	programWithFlags: programWithFlags
};

}();
var _elm_lang$virtual_dom$VirtualDom$programWithFlags = _elm_lang$virtual_dom$Native_VirtualDom.programWithFlags;
var _elm_lang$virtual_dom$VirtualDom$keyedNode = _elm_lang$virtual_dom$Native_VirtualDom.keyedNode;
var _elm_lang$virtual_dom$VirtualDom$lazy3 = _elm_lang$virtual_dom$Native_VirtualDom.lazy3;
var _elm_lang$virtual_dom$VirtualDom$lazy2 = _elm_lang$virtual_dom$Native_VirtualDom.lazy2;
var _elm_lang$virtual_dom$VirtualDom$lazy = _elm_lang$virtual_dom$Native_VirtualDom.lazy;
var _elm_lang$virtual_dom$VirtualDom$defaultOptions = {stopPropagation: false, preventDefault: false};
var _elm_lang$virtual_dom$VirtualDom$onWithOptions = _elm_lang$virtual_dom$Native_VirtualDom.on;
var _elm_lang$virtual_dom$VirtualDom$on = F2(
	function (eventName, decoder) {
		return A3(_elm_lang$virtual_dom$VirtualDom$onWithOptions, eventName, _elm_lang$virtual_dom$VirtualDom$defaultOptions, decoder);
	});
var _elm_lang$virtual_dom$VirtualDom$style = _elm_lang$virtual_dom$Native_VirtualDom.style;
var _elm_lang$virtual_dom$VirtualDom$attributeNS = _elm_lang$virtual_dom$Native_VirtualDom.attributeNS;
var _elm_lang$virtual_dom$VirtualDom$attribute = _elm_lang$virtual_dom$Native_VirtualDom.attribute;
var _elm_lang$virtual_dom$VirtualDom$property = _elm_lang$virtual_dom$Native_VirtualDom.property;
var _elm_lang$virtual_dom$VirtualDom$map = _elm_lang$virtual_dom$Native_VirtualDom.map;
var _elm_lang$virtual_dom$VirtualDom$text = _elm_lang$virtual_dom$Native_VirtualDom.text;
var _elm_lang$virtual_dom$VirtualDom$node = _elm_lang$virtual_dom$Native_VirtualDom.node;
var _elm_lang$virtual_dom$VirtualDom$Options = F2(
	function (a, b) {
		return {stopPropagation: a, preventDefault: b};
	});
var _elm_lang$virtual_dom$VirtualDom$Node = {ctor: 'Node'};
var _elm_lang$virtual_dom$VirtualDom$Property = {ctor: 'Property'};

var _elm_lang$html$Html$text = _elm_lang$virtual_dom$VirtualDom$text;
var _elm_lang$html$Html$node = _elm_lang$virtual_dom$VirtualDom$node;
var _elm_lang$html$Html$body = _elm_lang$html$Html$node('body');
var _elm_lang$html$Html$section = _elm_lang$html$Html$node('section');
var _elm_lang$html$Html$nav = _elm_lang$html$Html$node('nav');
var _elm_lang$html$Html$article = _elm_lang$html$Html$node('article');
var _elm_lang$html$Html$aside = _elm_lang$html$Html$node('aside');
var _elm_lang$html$Html$h1 = _elm_lang$html$Html$node('h1');
var _elm_lang$html$Html$h2 = _elm_lang$html$Html$node('h2');
var _elm_lang$html$Html$h3 = _elm_lang$html$Html$node('h3');
var _elm_lang$html$Html$h4 = _elm_lang$html$Html$node('h4');
var _elm_lang$html$Html$h5 = _elm_lang$html$Html$node('h5');
var _elm_lang$html$Html$h6 = _elm_lang$html$Html$node('h6');
var _elm_lang$html$Html$header = _elm_lang$html$Html$node('header');
var _elm_lang$html$Html$footer = _elm_lang$html$Html$node('footer');
var _elm_lang$html$Html$address = _elm_lang$html$Html$node('address');
var _elm_lang$html$Html$main$ = _elm_lang$html$Html$node('main');
var _elm_lang$html$Html$p = _elm_lang$html$Html$node('p');
var _elm_lang$html$Html$hr = _elm_lang$html$Html$node('hr');
var _elm_lang$html$Html$pre = _elm_lang$html$Html$node('pre');
var _elm_lang$html$Html$blockquote = _elm_lang$html$Html$node('blockquote');
var _elm_lang$html$Html$ol = _elm_lang$html$Html$node('ol');
var _elm_lang$html$Html$ul = _elm_lang$html$Html$node('ul');
var _elm_lang$html$Html$li = _elm_lang$html$Html$node('li');
var _elm_lang$html$Html$dl = _elm_lang$html$Html$node('dl');
var _elm_lang$html$Html$dt = _elm_lang$html$Html$node('dt');
var _elm_lang$html$Html$dd = _elm_lang$html$Html$node('dd');
var _elm_lang$html$Html$figure = _elm_lang$html$Html$node('figure');
var _elm_lang$html$Html$figcaption = _elm_lang$html$Html$node('figcaption');
var _elm_lang$html$Html$div = _elm_lang$html$Html$node('div');
var _elm_lang$html$Html$a = _elm_lang$html$Html$node('a');
var _elm_lang$html$Html$em = _elm_lang$html$Html$node('em');
var _elm_lang$html$Html$strong = _elm_lang$html$Html$node('strong');
var _elm_lang$html$Html$small = _elm_lang$html$Html$node('small');
var _elm_lang$html$Html$s = _elm_lang$html$Html$node('s');
var _elm_lang$html$Html$cite = _elm_lang$html$Html$node('cite');
var _elm_lang$html$Html$q = _elm_lang$html$Html$node('q');
var _elm_lang$html$Html$dfn = _elm_lang$html$Html$node('dfn');
var _elm_lang$html$Html$abbr = _elm_lang$html$Html$node('abbr');
var _elm_lang$html$Html$time = _elm_lang$html$Html$node('time');
var _elm_lang$html$Html$code = _elm_lang$html$Html$node('code');
var _elm_lang$html$Html$var = _elm_lang$html$Html$node('var');
var _elm_lang$html$Html$samp = _elm_lang$html$Html$node('samp');
var _elm_lang$html$Html$kbd = _elm_lang$html$Html$node('kbd');
var _elm_lang$html$Html$sub = _elm_lang$html$Html$node('sub');
var _elm_lang$html$Html$sup = _elm_lang$html$Html$node('sup');
var _elm_lang$html$Html$i = _elm_lang$html$Html$node('i');
var _elm_lang$html$Html$b = _elm_lang$html$Html$node('b');
var _elm_lang$html$Html$u = _elm_lang$html$Html$node('u');
var _elm_lang$html$Html$mark = _elm_lang$html$Html$node('mark');
var _elm_lang$html$Html$ruby = _elm_lang$html$Html$node('ruby');
var _elm_lang$html$Html$rt = _elm_lang$html$Html$node('rt');
var _elm_lang$html$Html$rp = _elm_lang$html$Html$node('rp');
var _elm_lang$html$Html$bdi = _elm_lang$html$Html$node('bdi');
var _elm_lang$html$Html$bdo = _elm_lang$html$Html$node('bdo');
var _elm_lang$html$Html$span = _elm_lang$html$Html$node('span');
var _elm_lang$html$Html$br = _elm_lang$html$Html$node('br');
var _elm_lang$html$Html$wbr = _elm_lang$html$Html$node('wbr');
var _elm_lang$html$Html$ins = _elm_lang$html$Html$node('ins');
var _elm_lang$html$Html$del = _elm_lang$html$Html$node('del');
var _elm_lang$html$Html$img = _elm_lang$html$Html$node('img');
var _elm_lang$html$Html$iframe = _elm_lang$html$Html$node('iframe');
var _elm_lang$html$Html$embed = _elm_lang$html$Html$node('embed');
var _elm_lang$html$Html$object = _elm_lang$html$Html$node('object');
var _elm_lang$html$Html$param = _elm_lang$html$Html$node('param');
var _elm_lang$html$Html$video = _elm_lang$html$Html$node('video');
var _elm_lang$html$Html$audio = _elm_lang$html$Html$node('audio');
var _elm_lang$html$Html$source = _elm_lang$html$Html$node('source');
var _elm_lang$html$Html$track = _elm_lang$html$Html$node('track');
var _elm_lang$html$Html$canvas = _elm_lang$html$Html$node('canvas');
var _elm_lang$html$Html$svg = _elm_lang$html$Html$node('svg');
var _elm_lang$html$Html$math = _elm_lang$html$Html$node('math');
var _elm_lang$html$Html$table = _elm_lang$html$Html$node('table');
var _elm_lang$html$Html$caption = _elm_lang$html$Html$node('caption');
var _elm_lang$html$Html$colgroup = _elm_lang$html$Html$node('colgroup');
var _elm_lang$html$Html$col = _elm_lang$html$Html$node('col');
var _elm_lang$html$Html$tbody = _elm_lang$html$Html$node('tbody');
var _elm_lang$html$Html$thead = _elm_lang$html$Html$node('thead');
var _elm_lang$html$Html$tfoot = _elm_lang$html$Html$node('tfoot');
var _elm_lang$html$Html$tr = _elm_lang$html$Html$node('tr');
var _elm_lang$html$Html$td = _elm_lang$html$Html$node('td');
var _elm_lang$html$Html$th = _elm_lang$html$Html$node('th');
var _elm_lang$html$Html$form = _elm_lang$html$Html$node('form');
var _elm_lang$html$Html$fieldset = _elm_lang$html$Html$node('fieldset');
var _elm_lang$html$Html$legend = _elm_lang$html$Html$node('legend');
var _elm_lang$html$Html$label = _elm_lang$html$Html$node('label');
var _elm_lang$html$Html$input = _elm_lang$html$Html$node('input');
var _elm_lang$html$Html$button = _elm_lang$html$Html$node('button');
var _elm_lang$html$Html$select = _elm_lang$html$Html$node('select');
var _elm_lang$html$Html$datalist = _elm_lang$html$Html$node('datalist');
var _elm_lang$html$Html$optgroup = _elm_lang$html$Html$node('optgroup');
var _elm_lang$html$Html$option = _elm_lang$html$Html$node('option');
var _elm_lang$html$Html$textarea = _elm_lang$html$Html$node('textarea');
var _elm_lang$html$Html$keygen = _elm_lang$html$Html$node('keygen');
var _elm_lang$html$Html$output = _elm_lang$html$Html$node('output');
var _elm_lang$html$Html$progress = _elm_lang$html$Html$node('progress');
var _elm_lang$html$Html$meter = _elm_lang$html$Html$node('meter');
var _elm_lang$html$Html$details = _elm_lang$html$Html$node('details');
var _elm_lang$html$Html$summary = _elm_lang$html$Html$node('summary');
var _elm_lang$html$Html$menuitem = _elm_lang$html$Html$node('menuitem');
var _elm_lang$html$Html$menu = _elm_lang$html$Html$node('menu');

var _elm_lang$html$Html_App$programWithFlags = _elm_lang$virtual_dom$VirtualDom$programWithFlags;
var _elm_lang$html$Html_App$program = function (app) {
	return _elm_lang$html$Html_App$programWithFlags(
		_elm_lang$core$Native_Utils.update(
			app,
			{
				init: function (_p0) {
					return app.init;
				}
			}));
};
var _elm_lang$html$Html_App$beginnerProgram = function (_p1) {
	var _p2 = _p1;
	return _elm_lang$html$Html_App$programWithFlags(
		{
			init: function (_p3) {
				return A2(
					_elm_lang$core$Platform_Cmd_ops['!'],
					_p2.model,
					_elm_lang$core$Native_List.fromArray(
						[]));
			},
			update: F2(
				function (msg, model) {
					return A2(
						_elm_lang$core$Platform_Cmd_ops['!'],
						A2(_p2.update, msg, model),
						_elm_lang$core$Native_List.fromArray(
							[]));
				}),
			view: _p2.view,
			subscriptions: function (_p4) {
				return _elm_lang$core$Platform_Sub$none;
			}
		});
};
var _elm_lang$html$Html_App$map = _elm_lang$virtual_dom$VirtualDom$map;

var _elm_lang$dom$Native_Dom = function() {

function on(node)
{
	return function(eventName, decoder, toTask)
	{
		return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback) {

			function performTask(event)
			{
				var result = A2(_elm_lang$core$Json_Decode$decodeValue, decoder, event);
				if (result.ctor === 'Ok')
				{
					_elm_lang$core$Native_Scheduler.rawSpawn(toTask(result._0));
				}
			}

			node.addEventListener(eventName, performTask);

			return function()
			{
				node.removeEventListener(eventName, performTask);
			};
		});
	};
}

return {
	onDocument: F3(on(document)),
	onWindow: F3(on(window))
};

}();

var _elm_lang$dom$Dom_LowLevel$onWindow = _elm_lang$dom$Native_Dom.onWindow;
var _elm_lang$dom$Dom_LowLevel$onDocument = _elm_lang$dom$Native_Dom.onDocument;

var _elm_lang$core$Process$kill = _elm_lang$core$Native_Scheduler.kill;
var _elm_lang$core$Process$sleep = _elm_lang$core$Native_Scheduler.sleep;
var _elm_lang$core$Process$spawn = _elm_lang$core$Native_Scheduler.spawn;

var _concourse$atc$Native_Scroll = function() {
  function toBottom() {
    return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback) {
      window.scrollTo(0, document.body.scrollHeight);
      callback(_elm_lang$core$Native_Scheduler.succeed(_elm_lang$core$Native_Utils.Tuple0));
    });
  }

  function scrollElement(id, delta) {
    return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback) {
      document.getElementById(id).scrollLeft -= delta;
      callback(_elm_lang$core$Native_Scheduler.succeed(_elm_lang$core$Native_Utils.Tuple0));
    });
  }

  function scrollIntoView(selector) {
    return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback) {
      document.querySelector(selector).scrollIntoView();
      callback(_elm_lang$core$Native_Scheduler.succeed(_elm_lang$core$Native_Utils.Tuple0));
    });
  }

  return {
    toBottom: toBottom,
    scrollElement: F2(scrollElement),
    scrollIntoView: scrollIntoView
  };
}();

var _concourse$atc$Scroll$onSelfMsg = F3(
	function (router, fb, state) {
		var _p0 = state;
		if (_p0.ctor === 'Nothing') {
			return _elm_lang$core$Task$succeed(_elm_lang$core$Maybe$Nothing);
		} else {
			var send = function (_p1) {
				var _p2 = _p1;
				return A2(
					_elm_lang$core$Platform$sendToApp,
					router,
					_p2._0(fb));
			};
			return A2(
				_elm_lang$core$Task$andThen,
				_elm_lang$core$Task$sequence(
					A2(_elm_lang$core$List$map, send, _p0._0.subs)),
				function (_p3) {
					return _elm_lang$core$Task$succeed(state);
				});
		}
	});
var _concourse$atc$Scroll$init = _elm_lang$core$Task$succeed(_elm_lang$core$Maybe$Nothing);
var _concourse$atc$Scroll$scrollIntoView = _concourse$atc$Native_Scroll.scrollIntoView;
var _concourse$atc$Scroll$scroll = _concourse$atc$Native_Scroll.scrollElement;
var _concourse$atc$Scroll$toBottom = _concourse$atc$Native_Scroll.toBottom(
	{ctor: '_Tuple0'});
var _concourse$atc$Scroll$decodeComparators = A4(
	_elm_lang$core$Json_Decode$object3,
	F3(
		function (v0, v1, v2) {
			return {ctor: '_Tuple3', _0: v0, _1: v1, _2: v2};
		}),
	A2(
		_elm_lang$core$Json_Decode$at,
		_elm_lang$core$Native_List.fromArray(
			['currentTarget', 'pageYOffset']),
		_elm_lang$core$Json_Decode$int),
	A2(
		_elm_lang$core$Json_Decode$at,
		_elm_lang$core$Native_List.fromArray(
			['target', 'documentElement', 'clientHeight']),
		_elm_lang$core$Json_Decode$int),
	A2(
		_elm_lang$core$Json_Decode$at,
		_elm_lang$core$Native_List.fromArray(
			['target', 'documentElement', 'scrollHeight']),
		_elm_lang$core$Json_Decode$int));
var _concourse$atc$Scroll$decodeFromBottom = A2(
	_elm_lang$core$Json_Decode$customDecoder,
	_concourse$atc$Scroll$decodeComparators,
	function (_p4) {
		var _p5 = _p4;
		var scrolledHeight = _p5._0 - _p5._1;
		return _elm_lang$core$Result$Ok(_p5._2 - scrolledHeight);
	});
var _concourse$atc$Scroll$onEffects = F3(
	function (router, subs, state) {
		var _p6 = state;
		if (_p6.ctor === 'Nothing') {
			var _p7 = subs;
			if (_p7.ctor === '[]') {
				return _elm_lang$core$Task$succeed(state);
			} else {
				return A2(
					_elm_lang$core$Task$andThen,
					_elm_lang$core$Process$spawn(
						A3(
							_elm_lang$dom$Dom_LowLevel$onWindow,
							'scroll',
							_concourse$atc$Scroll$decodeFromBottom,
							_elm_lang$core$Platform$sendToSelf(router))),
					function (watcher) {
						return _elm_lang$core$Task$succeed(
							_elm_lang$core$Maybe$Just(
								{subs: subs, watcher: watcher}));
					});
			}
		} else {
			var _p11 = _p6._0.watcher;
			var _p10 = _p6._0.subs;
			var _p8 = _p10;
			if (_p8.ctor === '[]') {
				return A2(
					_elm_lang$core$Task$andThen,
					_elm_lang$core$Process$kill(_p11),
					function (_p9) {
						return _elm_lang$core$Task$succeed(_elm_lang$core$Maybe$Nothing);
					});
			} else {
				return _elm_lang$core$Task$succeed(
					_elm_lang$core$Maybe$Just(
						{subs: _p10, watcher: _p11}));
			}
		}
	});
var _concourse$atc$Scroll$subscription = _elm_lang$core$Native_Platform.leaf('Scroll');
var _concourse$atc$Scroll$MySub = function (a) {
	return {ctor: 'MySub', _0: a};
};
var _concourse$atc$Scroll$fromBottom = function (tagger) {
	return _concourse$atc$Scroll$subscription(
		_concourse$atc$Scroll$MySub(tagger));
};
var _concourse$atc$Scroll$subMap = F2(
	function (func, _p12) {
		var _p13 = _p12;
		return _concourse$atc$Scroll$MySub(
			function (_p14) {
				return func(
					_p13._0(_p14));
			});
	});
_elm_lang$core$Native_Platform.effectManagers['Scroll'] = {pkg: 'concourse/atc', init: _concourse$atc$Scroll$init, onEffects: _concourse$atc$Scroll$onEffects, onSelfMsg: _concourse$atc$Scroll$onSelfMsg, tag: 'sub', subMap: _concourse$atc$Scroll$subMap};

var _concourse$atc$Autoscroll$Model = F3(
	function (a, b, c) {
		return {subModel: a, shouldScroll: b, scrollBehaviorFunc: c};
	});
var _concourse$atc$Autoscroll$NoScroll = {ctor: 'NoScroll'};
var _concourse$atc$Autoscroll$ScrollUntilCancelled = {ctor: 'ScrollUntilCancelled'};
var _concourse$atc$Autoscroll$Autoscroll = {ctor: 'Autoscroll'};
var _concourse$atc$Autoscroll$FromBottom = function (a) {
	return {ctor: 'FromBottom', _0: a};
};
var _concourse$atc$Autoscroll$ScrolledDown = {ctor: 'ScrolledDown'};
var _concourse$atc$Autoscroll$scrollToBottom = A3(
	_elm_lang$core$Task$perform,
	_elm_lang$core$Basics$always(_concourse$atc$Autoscroll$ScrolledDown),
	_elm_lang$core$Basics$always(_concourse$atc$Autoscroll$ScrolledDown),
	_concourse$atc$Scroll$toBottom);
var _concourse$atc$Autoscroll$ScrollDown = {ctor: 'ScrollDown'};
var _concourse$atc$Autoscroll$SubAction = function (a) {
	return {ctor: 'SubAction', _0: a};
};
var _concourse$atc$Autoscroll$init = F2(
	function (toScrollAction, _p0) {
		var _p1 = _p0;
		return {
			ctor: '_Tuple2',
			_0: A3(_concourse$atc$Autoscroll$Model, _p1._0, true, toScrollAction),
			_1: A2(_elm_lang$core$Platform_Cmd$map, _concourse$atc$Autoscroll$SubAction, _p1._1)
		};
	});
var _concourse$atc$Autoscroll$update = F3(
	function (subUpdate, action, model) {
		var _p2 = action;
		switch (_p2.ctor) {
			case 'SubAction':
				var _p3 = A2(subUpdate, _p2._0, model.subModel);
				var subModel = _p3._0;
				var subCmd = _p3._1;
				return {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Native_Utils.update(
						model,
						{subModel: subModel}),
					_1: A2(_elm_lang$core$Platform_Cmd$map, _concourse$atc$Autoscroll$SubAction, subCmd)
				};
			case 'ScrollDown':
				return {
					ctor: '_Tuple2',
					_0: model,
					_1: (model.shouldScroll && (!_elm_lang$core$Native_Utils.eq(
						model.scrollBehaviorFunc(model.subModel),
						_concourse$atc$Autoscroll$NoScroll))) ? _concourse$atc$Autoscroll$scrollToBottom : _elm_lang$core$Platform_Cmd$none
				};
			case 'ScrolledDown':
				return {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none};
			default:
				return {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Native_Utils.update(
						model,
						{
							shouldScroll: function () {
								var _p4 = model.scrollBehaviorFunc(model.subModel);
								if (_p4.ctor === 'Autoscroll') {
									return _elm_lang$core$Native_Utils.cmp(_p2._0, 16) < 0;
								} else {
									return false;
								}
							}()
						}),
					_1: _elm_lang$core$Platform_Cmd$none
				};
		}
	});
var _concourse$atc$Autoscroll$view = F2(
	function (subView, model) {
		return A2(
			_elm_lang$html$Html_App$map,
			_concourse$atc$Autoscroll$SubAction,
			subView(model.subModel));
	});

//import Result //

var _elm_lang$core$Native_Date = function() {

function fromString(str)
{
	var date = new Date(str);
	return isNaN(date.getTime())
		? _elm_lang$core$Result$Err('Unable to parse \'' + str + '\' as a date. Dates must be in the ISO 8601 format.')
		: _elm_lang$core$Result$Ok(date);
}

var dayTable = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
var monthTable =
	['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
	 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];


return {
	fromString: fromString,
	year: function(d) { return d.getFullYear(); },
	month: function(d) { return { ctor: monthTable[d.getMonth()] }; },
	day: function(d) { return d.getDate(); },
	hour: function(d) { return d.getHours(); },
	minute: function(d) { return d.getMinutes(); },
	second: function(d) { return d.getSeconds(); },
	millisecond: function(d) { return d.getMilliseconds(); },
	toTime: function(d) { return d.getTime(); },
	fromTime: function(t) { return new Date(t); },
	dayOfWeek: function(d) { return { ctor: dayTable[d.getDay()] }; }
};

}();
var _elm_lang$core$Date$millisecond = _elm_lang$core$Native_Date.millisecond;
var _elm_lang$core$Date$second = _elm_lang$core$Native_Date.second;
var _elm_lang$core$Date$minute = _elm_lang$core$Native_Date.minute;
var _elm_lang$core$Date$hour = _elm_lang$core$Native_Date.hour;
var _elm_lang$core$Date$dayOfWeek = _elm_lang$core$Native_Date.dayOfWeek;
var _elm_lang$core$Date$day = _elm_lang$core$Native_Date.day;
var _elm_lang$core$Date$month = _elm_lang$core$Native_Date.month;
var _elm_lang$core$Date$year = _elm_lang$core$Native_Date.year;
var _elm_lang$core$Date$fromTime = _elm_lang$core$Native_Date.fromTime;
var _elm_lang$core$Date$toTime = _elm_lang$core$Native_Date.toTime;
var _elm_lang$core$Date$fromString = _elm_lang$core$Native_Date.fromString;
var _elm_lang$core$Date$now = A2(_elm_lang$core$Task$map, _elm_lang$core$Date$fromTime, _elm_lang$core$Time$now);
var _elm_lang$core$Date$Date = {ctor: 'Date'};
var _elm_lang$core$Date$Sun = {ctor: 'Sun'};
var _elm_lang$core$Date$Sat = {ctor: 'Sat'};
var _elm_lang$core$Date$Fri = {ctor: 'Fri'};
var _elm_lang$core$Date$Thu = {ctor: 'Thu'};
var _elm_lang$core$Date$Wed = {ctor: 'Wed'};
var _elm_lang$core$Date$Tue = {ctor: 'Tue'};
var _elm_lang$core$Date$Mon = {ctor: 'Mon'};
var _elm_lang$core$Date$Dec = {ctor: 'Dec'};
var _elm_lang$core$Date$Nov = {ctor: 'Nov'};
var _elm_lang$core$Date$Oct = {ctor: 'Oct'};
var _elm_lang$core$Date$Sep = {ctor: 'Sep'};
var _elm_lang$core$Date$Aug = {ctor: 'Aug'};
var _elm_lang$core$Date$Jul = {ctor: 'Jul'};
var _elm_lang$core$Date$Jun = {ctor: 'Jun'};
var _elm_lang$core$Date$May = {ctor: 'May'};
var _elm_lang$core$Date$Apr = {ctor: 'Apr'};
var _elm_lang$core$Date$Mar = {ctor: 'Mar'};
var _elm_lang$core$Date$Feb = {ctor: 'Feb'};
var _elm_lang$core$Date$Jan = {ctor: 'Jan'};

//import Maybe, Native.List //

var _elm_lang$core$Native_Regex = function() {

function escape(str)
{
	return str.replace(/[-\/\\^$*+?.()|[\]{}]/g, '\\$&');
}
function caseInsensitive(re)
{
	return new RegExp(re.source, 'gi');
}
function regex(raw)
{
	return new RegExp(raw, 'g');
}

function contains(re, string)
{
	return string.match(re) !== null;
}

function find(n, re, str)
{
	n = n.ctor === 'All' ? Infinity : n._0;
	var out = [];
	var number = 0;
	var string = str;
	var lastIndex = re.lastIndex;
	var prevLastIndex = -1;
	var result;
	while (number++ < n && (result = re.exec(string)))
	{
		if (prevLastIndex === re.lastIndex) break;
		var i = result.length - 1;
		var subs = new Array(i);
		while (i > 0)
		{
			var submatch = result[i];
			subs[--i] = submatch === undefined
				? _elm_lang$core$Maybe$Nothing
				: _elm_lang$core$Maybe$Just(submatch);
		}
		out.push({
			match: result[0],
			submatches: _elm_lang$core$Native_List.fromArray(subs),
			index: result.index,
			number: number
		});
		prevLastIndex = re.lastIndex;
	}
	re.lastIndex = lastIndex;
	return _elm_lang$core$Native_List.fromArray(out);
}

function replace(n, re, replacer, string)
{
	n = n.ctor === 'All' ? Infinity : n._0;
	var count = 0;
	function jsReplacer(match)
	{
		if (count++ >= n)
		{
			return match;
		}
		var i = arguments.length - 3;
		var submatches = new Array(i);
		while (i > 0)
		{
			var submatch = arguments[i];
			submatches[--i] = submatch === undefined
				? _elm_lang$core$Maybe$Nothing
				: _elm_lang$core$Maybe$Just(submatch);
		}
		return replacer({
			match: match,
			submatches: _elm_lang$core$Native_List.fromArray(submatches),
			index: arguments[i - 1],
			number: count
		});
	}
	return string.replace(re, jsReplacer);
}

function split(n, re, str)
{
	n = n.ctor === 'All' ? Infinity : n._0;
	if (n === Infinity)
	{
		return _elm_lang$core$Native_List.fromArray(str.split(re));
	}
	var string = str;
	var result;
	var out = [];
	var start = re.lastIndex;
	while (n--)
	{
		if (!(result = re.exec(string))) break;
		out.push(string.slice(start, result.index));
		start = re.lastIndex;
	}
	out.push(string.slice(start));
	return _elm_lang$core$Native_List.fromArray(out);
}

return {
	regex: regex,
	caseInsensitive: caseInsensitive,
	escape: escape,

	contains: F2(contains),
	find: F3(find),
	replace: F4(replace),
	split: F3(split)
};

}();

var _elm_lang$core$Regex$split = _elm_lang$core$Native_Regex.split;
var _elm_lang$core$Regex$replace = _elm_lang$core$Native_Regex.replace;
var _elm_lang$core$Regex$find = _elm_lang$core$Native_Regex.find;
var _elm_lang$core$Regex$contains = _elm_lang$core$Native_Regex.contains;
var _elm_lang$core$Regex$caseInsensitive = _elm_lang$core$Native_Regex.caseInsensitive;
var _elm_lang$core$Regex$regex = _elm_lang$core$Native_Regex.regex;
var _elm_lang$core$Regex$escape = _elm_lang$core$Native_Regex.escape;
var _elm_lang$core$Regex$Match = F4(
	function (a, b, c, d) {
		return {match: a, submatches: b, index: c, number: d};
	});
var _elm_lang$core$Regex$Regex = {ctor: 'Regex'};
var _elm_lang$core$Regex$AtMost = function (a) {
	return {ctor: 'AtMost', _0: a};
};
var _elm_lang$core$Regex$All = {ctor: 'All'};

var _mgold$elm_date_format$Date_Format$padWith = function (c) {
	return function (_p0) {
		return A3(
			_elm_lang$core$String$padLeft,
			2,
			c,
			_elm_lang$core$Basics$toString(_p0));
	};
};
var _mgold$elm_date_format$Date_Format$zero2twelve = function (n) {
	return _elm_lang$core$Native_Utils.eq(n, 0) ? 12 : n;
};
var _mgold$elm_date_format$Date_Format$mod12 = function (h) {
	return A2(_elm_lang$core$Basics_ops['%'], h, 12);
};
var _mgold$elm_date_format$Date_Format$fullDayOfWeek = function (dow) {
	var _p1 = dow;
	switch (_p1.ctor) {
		case 'Mon':
			return 'Monday';
		case 'Tue':
			return 'Tuesday';
		case 'Wed':
			return 'Wednesday';
		case 'Thu':
			return 'Thursday';
		case 'Fri':
			return 'Friday';
		case 'Sat':
			return 'Saturday';
		default:
			return 'Sunday';
	}
};
var _mgold$elm_date_format$Date_Format$monthToFullName = function (m) {
	var _p2 = m;
	switch (_p2.ctor) {
		case 'Jan':
			return 'January';
		case 'Feb':
			return 'February';
		case 'Mar':
			return 'March';
		case 'Apr':
			return 'April';
		case 'May':
			return 'May';
		case 'Jun':
			return 'June';
		case 'Jul':
			return 'July';
		case 'Aug':
			return 'August';
		case 'Sep':
			return 'September';
		case 'Oct':
			return 'October';
		case 'Nov':
			return 'November';
		default:
			return 'December';
	}
};
var _mgold$elm_date_format$Date_Format$monthToInt = function (m) {
	var _p3 = m;
	switch (_p3.ctor) {
		case 'Jan':
			return 1;
		case 'Feb':
			return 2;
		case 'Mar':
			return 3;
		case 'Apr':
			return 4;
		case 'May':
			return 5;
		case 'Jun':
			return 6;
		case 'Jul':
			return 7;
		case 'Aug':
			return 8;
		case 'Sep':
			return 9;
		case 'Oct':
			return 10;
		case 'Nov':
			return 11;
		default:
			return 12;
	}
};
var _mgold$elm_date_format$Date_Format$formatToken = F2(
	function (d, m) {
		var symbol = function () {
			var _p4 = m.submatches;
			if (((_p4.ctor === '::') && (_p4._0.ctor === 'Just')) && (_p4._1.ctor === '[]')) {
				return _p4._0._0;
			} else {
				return ' ';
			}
		}();
		var _p5 = symbol;
		switch (_p5) {
			case '%':
				return '%';
			case 'Y':
				return _elm_lang$core$Basics$toString(
					_elm_lang$core$Date$year(d));
			case 'm':
				return A3(
					_elm_lang$core$String$padLeft,
					2,
					_elm_lang$core$Native_Utils.chr('0'),
					_elm_lang$core$Basics$toString(
						_mgold$elm_date_format$Date_Format$monthToInt(
							_elm_lang$core$Date$month(d))));
			case 'B':
				return _mgold$elm_date_format$Date_Format$monthToFullName(
					_elm_lang$core$Date$month(d));
			case 'b':
				return _elm_lang$core$Basics$toString(
					_elm_lang$core$Date$month(d));
			case 'd':
				return A2(
					_mgold$elm_date_format$Date_Format$padWith,
					_elm_lang$core$Native_Utils.chr('0'),
					_elm_lang$core$Date$day(d));
			case 'e':
				return A2(
					_mgold$elm_date_format$Date_Format$padWith,
					_elm_lang$core$Native_Utils.chr(' '),
					_elm_lang$core$Date$day(d));
			case 'a':
				return _elm_lang$core$Basics$toString(
					_elm_lang$core$Date$dayOfWeek(d));
			case 'A':
				return _mgold$elm_date_format$Date_Format$fullDayOfWeek(
					_elm_lang$core$Date$dayOfWeek(d));
			case 'H':
				return A2(
					_mgold$elm_date_format$Date_Format$padWith,
					_elm_lang$core$Native_Utils.chr('0'),
					_elm_lang$core$Date$hour(d));
			case 'k':
				return A2(
					_mgold$elm_date_format$Date_Format$padWith,
					_elm_lang$core$Native_Utils.chr(' '),
					_elm_lang$core$Date$hour(d));
			case 'I':
				return A2(
					_mgold$elm_date_format$Date_Format$padWith,
					_elm_lang$core$Native_Utils.chr('0'),
					_mgold$elm_date_format$Date_Format$zero2twelve(
						_mgold$elm_date_format$Date_Format$mod12(
							_elm_lang$core$Date$hour(d))));
			case 'l':
				return A2(
					_mgold$elm_date_format$Date_Format$padWith,
					_elm_lang$core$Native_Utils.chr(' '),
					_mgold$elm_date_format$Date_Format$zero2twelve(
						_mgold$elm_date_format$Date_Format$mod12(
							_elm_lang$core$Date$hour(d))));
			case 'p':
				return (_elm_lang$core$Native_Utils.cmp(
					_elm_lang$core$Date$hour(d),
					13) < 0) ? 'AM' : 'PM';
			case 'P':
				return (_elm_lang$core$Native_Utils.cmp(
					_elm_lang$core$Date$hour(d),
					13) < 0) ? 'am' : 'pm';
			case 'M':
				return A2(
					_mgold$elm_date_format$Date_Format$padWith,
					_elm_lang$core$Native_Utils.chr('0'),
					_elm_lang$core$Date$minute(d));
			case 'S':
				return A2(
					_mgold$elm_date_format$Date_Format$padWith,
					_elm_lang$core$Native_Utils.chr('0'),
					_elm_lang$core$Date$second(d));
			default:
				return '';
		}
	});
var _mgold$elm_date_format$Date_Format$re = _elm_lang$core$Regex$regex('%(%|Y|m|B|b|d|e|a|A|H|k|I|l|p|P|M|S)');
var _mgold$elm_date_format$Date_Format$format = F2(
	function (s, d) {
		return A4(
			_elm_lang$core$Regex$replace,
			_elm_lang$core$Regex$All,
			_mgold$elm_date_format$Date_Format$re,
			_mgold$elm_date_format$Date_Format$formatToken(d),
			s);
	});
var _mgold$elm_date_format$Date_Format$formatISO8601 = _mgold$elm_date_format$Date_Format$format('%Y-%m-%dT%H:%M:%SZ');

var _elm_lang$html$Html_Attributes$attribute = _elm_lang$virtual_dom$VirtualDom$attribute;
var _elm_lang$html$Html_Attributes$contextmenu = function (value) {
	return A2(_elm_lang$html$Html_Attributes$attribute, 'contextmenu', value);
};
var _elm_lang$html$Html_Attributes$draggable = function (value) {
	return A2(_elm_lang$html$Html_Attributes$attribute, 'draggable', value);
};
var _elm_lang$html$Html_Attributes$list = function (value) {
	return A2(_elm_lang$html$Html_Attributes$attribute, 'list', value);
};
var _elm_lang$html$Html_Attributes$maxlength = function (n) {
	return A2(
		_elm_lang$html$Html_Attributes$attribute,
		'maxlength',
		_elm_lang$core$Basics$toString(n));
};
var _elm_lang$html$Html_Attributes$datetime = function (value) {
	return A2(_elm_lang$html$Html_Attributes$attribute, 'datetime', value);
};
var _elm_lang$html$Html_Attributes$pubdate = function (value) {
	return A2(_elm_lang$html$Html_Attributes$attribute, 'pubdate', value);
};
var _elm_lang$html$Html_Attributes$colspan = function (n) {
	return A2(
		_elm_lang$html$Html_Attributes$attribute,
		'colspan',
		_elm_lang$core$Basics$toString(n));
};
var _elm_lang$html$Html_Attributes$rowspan = function (n) {
	return A2(
		_elm_lang$html$Html_Attributes$attribute,
		'rowspan',
		_elm_lang$core$Basics$toString(n));
};
var _elm_lang$html$Html_Attributes$property = _elm_lang$virtual_dom$VirtualDom$property;
var _elm_lang$html$Html_Attributes$stringProperty = F2(
	function (name, string) {
		return A2(
			_elm_lang$html$Html_Attributes$property,
			name,
			_elm_lang$core$Json_Encode$string(string));
	});
var _elm_lang$html$Html_Attributes$class = function (name) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'className', name);
};
var _elm_lang$html$Html_Attributes$id = function (name) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'id', name);
};
var _elm_lang$html$Html_Attributes$title = function (name) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'title', name);
};
var _elm_lang$html$Html_Attributes$accesskey = function ($char) {
	return A2(
		_elm_lang$html$Html_Attributes$stringProperty,
		'accessKey',
		_elm_lang$core$String$fromChar($char));
};
var _elm_lang$html$Html_Attributes$dir = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'dir', value);
};
var _elm_lang$html$Html_Attributes$dropzone = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'dropzone', value);
};
var _elm_lang$html$Html_Attributes$itemprop = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'itemprop', value);
};
var _elm_lang$html$Html_Attributes$lang = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'lang', value);
};
var _elm_lang$html$Html_Attributes$tabindex = function (n) {
	return A2(
		_elm_lang$html$Html_Attributes$stringProperty,
		'tabIndex',
		_elm_lang$core$Basics$toString(n));
};
var _elm_lang$html$Html_Attributes$charset = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'charset', value);
};
var _elm_lang$html$Html_Attributes$content = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'content', value);
};
var _elm_lang$html$Html_Attributes$httpEquiv = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'httpEquiv', value);
};
var _elm_lang$html$Html_Attributes$language = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'language', value);
};
var _elm_lang$html$Html_Attributes$src = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'src', value);
};
var _elm_lang$html$Html_Attributes$height = function (value) {
	return A2(
		_elm_lang$html$Html_Attributes$stringProperty,
		'height',
		_elm_lang$core$Basics$toString(value));
};
var _elm_lang$html$Html_Attributes$width = function (value) {
	return A2(
		_elm_lang$html$Html_Attributes$stringProperty,
		'width',
		_elm_lang$core$Basics$toString(value));
};
var _elm_lang$html$Html_Attributes$alt = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'alt', value);
};
var _elm_lang$html$Html_Attributes$preload = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'preload', value);
};
var _elm_lang$html$Html_Attributes$poster = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'poster', value);
};
var _elm_lang$html$Html_Attributes$kind = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'kind', value);
};
var _elm_lang$html$Html_Attributes$srclang = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'srclang', value);
};
var _elm_lang$html$Html_Attributes$sandbox = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'sandbox', value);
};
var _elm_lang$html$Html_Attributes$srcdoc = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'srcdoc', value);
};
var _elm_lang$html$Html_Attributes$type$ = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'type', value);
};
var _elm_lang$html$Html_Attributes$value = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'value', value);
};
var _elm_lang$html$Html_Attributes$defaultValue = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'defaultValue', value);
};
var _elm_lang$html$Html_Attributes$placeholder = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'placeholder', value);
};
var _elm_lang$html$Html_Attributes$accept = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'accept', value);
};
var _elm_lang$html$Html_Attributes$acceptCharset = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'acceptCharset', value);
};
var _elm_lang$html$Html_Attributes$action = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'action', value);
};
var _elm_lang$html$Html_Attributes$autocomplete = function (bool) {
	return A2(
		_elm_lang$html$Html_Attributes$stringProperty,
		'autocomplete',
		bool ? 'on' : 'off');
};
var _elm_lang$html$Html_Attributes$autosave = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'autosave', value);
};
var _elm_lang$html$Html_Attributes$enctype = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'enctype', value);
};
var _elm_lang$html$Html_Attributes$formaction = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'formAction', value);
};
var _elm_lang$html$Html_Attributes$minlength = function (n) {
	return A2(
		_elm_lang$html$Html_Attributes$stringProperty,
		'minLength',
		_elm_lang$core$Basics$toString(n));
};
var _elm_lang$html$Html_Attributes$method = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'method', value);
};
var _elm_lang$html$Html_Attributes$name = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'name', value);
};
var _elm_lang$html$Html_Attributes$pattern = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'pattern', value);
};
var _elm_lang$html$Html_Attributes$size = function (n) {
	return A2(
		_elm_lang$html$Html_Attributes$stringProperty,
		'size',
		_elm_lang$core$Basics$toString(n));
};
var _elm_lang$html$Html_Attributes$for = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'htmlFor', value);
};
var _elm_lang$html$Html_Attributes$form = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'form', value);
};
var _elm_lang$html$Html_Attributes$max = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'max', value);
};
var _elm_lang$html$Html_Attributes$min = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'min', value);
};
var _elm_lang$html$Html_Attributes$step = function (n) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'step', n);
};
var _elm_lang$html$Html_Attributes$cols = function (n) {
	return A2(
		_elm_lang$html$Html_Attributes$stringProperty,
		'cols',
		_elm_lang$core$Basics$toString(n));
};
var _elm_lang$html$Html_Attributes$rows = function (n) {
	return A2(
		_elm_lang$html$Html_Attributes$stringProperty,
		'rows',
		_elm_lang$core$Basics$toString(n));
};
var _elm_lang$html$Html_Attributes$wrap = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'wrap', value);
};
var _elm_lang$html$Html_Attributes$usemap = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'useMap', value);
};
var _elm_lang$html$Html_Attributes$shape = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'shape', value);
};
var _elm_lang$html$Html_Attributes$coords = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'coords', value);
};
var _elm_lang$html$Html_Attributes$challenge = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'challenge', value);
};
var _elm_lang$html$Html_Attributes$keytype = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'keytype', value);
};
var _elm_lang$html$Html_Attributes$align = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'align', value);
};
var _elm_lang$html$Html_Attributes$cite = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'cite', value);
};
var _elm_lang$html$Html_Attributes$href = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'href', value);
};
var _elm_lang$html$Html_Attributes$target = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'target', value);
};
var _elm_lang$html$Html_Attributes$downloadAs = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'download', value);
};
var _elm_lang$html$Html_Attributes$hreflang = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'hreflang', value);
};
var _elm_lang$html$Html_Attributes$media = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'media', value);
};
var _elm_lang$html$Html_Attributes$ping = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'ping', value);
};
var _elm_lang$html$Html_Attributes$rel = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'rel', value);
};
var _elm_lang$html$Html_Attributes$start = function (n) {
	return A2(
		_elm_lang$html$Html_Attributes$stringProperty,
		'start',
		_elm_lang$core$Basics$toString(n));
};
var _elm_lang$html$Html_Attributes$headers = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'headers', value);
};
var _elm_lang$html$Html_Attributes$scope = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'scope', value);
};
var _elm_lang$html$Html_Attributes$manifest = function (value) {
	return A2(_elm_lang$html$Html_Attributes$stringProperty, 'manifest', value);
};
var _elm_lang$html$Html_Attributes$boolProperty = F2(
	function (name, bool) {
		return A2(
			_elm_lang$html$Html_Attributes$property,
			name,
			_elm_lang$core$Json_Encode$bool(bool));
	});
var _elm_lang$html$Html_Attributes$hidden = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'hidden', bool);
};
var _elm_lang$html$Html_Attributes$contenteditable = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'contentEditable', bool);
};
var _elm_lang$html$Html_Attributes$spellcheck = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'spellcheck', bool);
};
var _elm_lang$html$Html_Attributes$async = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'async', bool);
};
var _elm_lang$html$Html_Attributes$defer = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'defer', bool);
};
var _elm_lang$html$Html_Attributes$scoped = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'scoped', bool);
};
var _elm_lang$html$Html_Attributes$autoplay = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'autoplay', bool);
};
var _elm_lang$html$Html_Attributes$controls = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'controls', bool);
};
var _elm_lang$html$Html_Attributes$loop = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'loop', bool);
};
var _elm_lang$html$Html_Attributes$default = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'default', bool);
};
var _elm_lang$html$Html_Attributes$seamless = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'seamless', bool);
};
var _elm_lang$html$Html_Attributes$checked = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'checked', bool);
};
var _elm_lang$html$Html_Attributes$selected = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'selected', bool);
};
var _elm_lang$html$Html_Attributes$autofocus = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'autofocus', bool);
};
var _elm_lang$html$Html_Attributes$disabled = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'disabled', bool);
};
var _elm_lang$html$Html_Attributes$multiple = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'multiple', bool);
};
var _elm_lang$html$Html_Attributes$novalidate = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'noValidate', bool);
};
var _elm_lang$html$Html_Attributes$readonly = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'readOnly', bool);
};
var _elm_lang$html$Html_Attributes$required = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'required', bool);
};
var _elm_lang$html$Html_Attributes$ismap = function (value) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'isMap', value);
};
var _elm_lang$html$Html_Attributes$download = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'download', bool);
};
var _elm_lang$html$Html_Attributes$reversed = function (bool) {
	return A2(_elm_lang$html$Html_Attributes$boolProperty, 'reversed', bool);
};
var _elm_lang$html$Html_Attributes$classList = function (list) {
	return _elm_lang$html$Html_Attributes$class(
		A2(
			_elm_lang$core$String$join,
			' ',
			A2(
				_elm_lang$core$List$map,
				_elm_lang$core$Basics$fst,
				A2(_elm_lang$core$List$filter, _elm_lang$core$Basics$snd, list))));
};
var _elm_lang$html$Html_Attributes$style = _elm_lang$virtual_dom$VirtualDom$style;

var _elm_lang$html$Html_Events$keyCode = A2(_elm_lang$core$Json_Decode_ops[':='], 'keyCode', _elm_lang$core$Json_Decode$int);
var _elm_lang$html$Html_Events$targetChecked = A2(
	_elm_lang$core$Json_Decode$at,
	_elm_lang$core$Native_List.fromArray(
		['target', 'checked']),
	_elm_lang$core$Json_Decode$bool);
var _elm_lang$html$Html_Events$targetValue = A2(
	_elm_lang$core$Json_Decode$at,
	_elm_lang$core$Native_List.fromArray(
		['target', 'value']),
	_elm_lang$core$Json_Decode$string);
var _elm_lang$html$Html_Events$defaultOptions = _elm_lang$virtual_dom$VirtualDom$defaultOptions;
var _elm_lang$html$Html_Events$onWithOptions = _elm_lang$virtual_dom$VirtualDom$onWithOptions;
var _elm_lang$html$Html_Events$on = _elm_lang$virtual_dom$VirtualDom$on;
var _elm_lang$html$Html_Events$onFocus = function (msg) {
	return A2(
		_elm_lang$html$Html_Events$on,
		'focus',
		_elm_lang$core$Json_Decode$succeed(msg));
};
var _elm_lang$html$Html_Events$onBlur = function (msg) {
	return A2(
		_elm_lang$html$Html_Events$on,
		'blur',
		_elm_lang$core$Json_Decode$succeed(msg));
};
var _elm_lang$html$Html_Events$onSubmitOptions = _elm_lang$core$Native_Utils.update(
	_elm_lang$html$Html_Events$defaultOptions,
	{preventDefault: true});
var _elm_lang$html$Html_Events$onSubmit = function (msg) {
	return A3(
		_elm_lang$html$Html_Events$onWithOptions,
		'submit',
		_elm_lang$html$Html_Events$onSubmitOptions,
		_elm_lang$core$Json_Decode$succeed(msg));
};
var _elm_lang$html$Html_Events$onCheck = function (tagger) {
	return A2(
		_elm_lang$html$Html_Events$on,
		'change',
		A2(_elm_lang$core$Json_Decode$map, tagger, _elm_lang$html$Html_Events$targetChecked));
};
var _elm_lang$html$Html_Events$onInput = function (tagger) {
	return A2(
		_elm_lang$html$Html_Events$on,
		'input',
		A2(_elm_lang$core$Json_Decode$map, tagger, _elm_lang$html$Html_Events$targetValue));
};
var _elm_lang$html$Html_Events$onMouseOut = function (msg) {
	return A2(
		_elm_lang$html$Html_Events$on,
		'mouseout',
		_elm_lang$core$Json_Decode$succeed(msg));
};
var _elm_lang$html$Html_Events$onMouseOver = function (msg) {
	return A2(
		_elm_lang$html$Html_Events$on,
		'mouseover',
		_elm_lang$core$Json_Decode$succeed(msg));
};
var _elm_lang$html$Html_Events$onMouseLeave = function (msg) {
	return A2(
		_elm_lang$html$Html_Events$on,
		'mouseleave',
		_elm_lang$core$Json_Decode$succeed(msg));
};
var _elm_lang$html$Html_Events$onMouseEnter = function (msg) {
	return A2(
		_elm_lang$html$Html_Events$on,
		'mouseenter',
		_elm_lang$core$Json_Decode$succeed(msg));
};
var _elm_lang$html$Html_Events$onMouseUp = function (msg) {
	return A2(
		_elm_lang$html$Html_Events$on,
		'mouseup',
		_elm_lang$core$Json_Decode$succeed(msg));
};
var _elm_lang$html$Html_Events$onMouseDown = function (msg) {
	return A2(
		_elm_lang$html$Html_Events$on,
		'mousedown',
		_elm_lang$core$Json_Decode$succeed(msg));
};
var _elm_lang$html$Html_Events$onDoubleClick = function (msg) {
	return A2(
		_elm_lang$html$Html_Events$on,
		'dblclick',
		_elm_lang$core$Json_Decode$succeed(msg));
};
var _elm_lang$html$Html_Events$onClick = function (msg) {
	return A2(
		_elm_lang$html$Html_Events$on,
		'click',
		_elm_lang$core$Json_Decode$succeed(msg));
};
var _elm_lang$html$Html_Events$Options = F2(
	function (a, b) {
		return {stopPropagation: a, preventDefault: b};
	});

var _elm_lang$html$Html_Lazy$lazy3 = _elm_lang$virtual_dom$VirtualDom$lazy3;
var _elm_lang$html$Html_Lazy$lazy2 = _elm_lang$virtual_dom$VirtualDom$lazy2;
var _elm_lang$html$Html_Lazy$lazy = _elm_lang$virtual_dom$VirtualDom$lazy;

//import Dict, List, Maybe, Native.Scheduler //

var _evancz$elm_http$Native_Http = function() {

function send(settings, request)
{
	return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback) {
		var req = new XMLHttpRequest();

		// start
		if (settings.onStart.ctor === 'Just')
		{
			req.addEventListener('loadStart', function() {
				var task = settings.onStart._0;
				_elm_lang$core$Native_Scheduler.rawSpawn(task);
			});
		}

		// progress
		if (settings.onProgress.ctor === 'Just')
		{
			req.addEventListener('progress', function(event) {
				var progress = !event.lengthComputable
					? _elm_lang$core$Maybe$Nothing
					: _elm_lang$core$Maybe$Just({
						loaded: event.loaded,
						total: event.total
					});
				var task = settings.onProgress._0(progress);
				_elm_lang$core$Native_Scheduler.rawSpawn(task);
			});
		}

		// end
		req.addEventListener('error', function() {
			return callback(_elm_lang$core$Native_Scheduler.fail({ ctor: 'RawNetworkError' }));
		});

		req.addEventListener('timeout', function() {
			return callback(_elm_lang$core$Native_Scheduler.fail({ ctor: 'RawTimeout' }));
		});

		req.addEventListener('load', function() {
			return callback(_elm_lang$core$Native_Scheduler.succeed(toResponse(req)));
		});

		req.open(request.verb, request.url, true);

		// set all the headers
		function setHeader(pair) {
			req.setRequestHeader(pair._0, pair._1);
		}
		A2(_elm_lang$core$List$map, setHeader, request.headers);

		// set the timeout
		req.timeout = settings.timeout;

		// enable this withCredentials thing
		req.withCredentials = settings.withCredentials;

		// ask for a specific MIME type for the response
		if (settings.desiredResponseType.ctor === 'Just')
		{
			req.overrideMimeType(settings.desiredResponseType._0);
		}

		// actuall send the request
		if(request.body.ctor === "BodyFormData")
		{
			req.send(request.body.formData)
		}
		else
		{
			req.send(request.body._0);
		}

		return function() {
			req.abort();
		};
	});
}


// deal with responses

function toResponse(req)
{
	var tag = req.responseType === 'blob' ? 'Blob' : 'Text'
	var response = tag === 'Blob' ? req.response : req.responseText;
	return {
		status: req.status,
		statusText: req.statusText,
		headers: parseHeaders(req.getAllResponseHeaders()),
		url: req.responseURL,
		value: { ctor: tag, _0: response }
	};
}


function parseHeaders(rawHeaders)
{
	var headers = _elm_lang$core$Dict$empty;

	if (!rawHeaders)
	{
		return headers;
	}

	var headerPairs = rawHeaders.split('\u000d\u000a');
	for (var i = headerPairs.length; i--; )
	{
		var headerPair = headerPairs[i];
		var index = headerPair.indexOf('\u003a\u0020');
		if (index > 0)
		{
			var key = headerPair.substring(0, index);
			var value = headerPair.substring(index + 2);

			headers = A3(_elm_lang$core$Dict$update, key, function(oldValue) {
				if (oldValue.ctor === 'Just')
				{
					return _elm_lang$core$Maybe$Just(value + ', ' + oldValue._0);
				}
				return _elm_lang$core$Maybe$Just(value);
			}, headers);
		}
	}

	return headers;
}


function multipart(dataList)
{
	var formData = new FormData();

	while (dataList.ctor !== '[]')
	{
		var data = dataList._0;
		if (data.ctor === 'StringData')
		{
			formData.append(data._0, data._1);
		}
		else
		{
			var fileName = data._1.ctor === 'Nothing'
				? undefined
				: data._1._0;
			formData.append(data._0, data._2, fileName);
		}
		dataList = dataList._1;
	}

	return { ctor: 'BodyFormData', formData: formData };
}


function uriEncode(string)
{
	return encodeURIComponent(string);
}

function uriDecode(string)
{
	return decodeURIComponent(string);
}

return {
	send: F2(send),
	multipart: multipart,
	uriEncode: uriEncode,
	uriDecode: uriDecode
};

}();

var _evancz$elm_http$Http$send = _evancz$elm_http$Native_Http.send;
var _evancz$elm_http$Http$defaultSettings = {timeout: 0, onStart: _elm_lang$core$Maybe$Nothing, onProgress: _elm_lang$core$Maybe$Nothing, desiredResponseType: _elm_lang$core$Maybe$Nothing, withCredentials: false};
var _evancz$elm_http$Http$multipart = _evancz$elm_http$Native_Http.multipart;
var _evancz$elm_http$Http$uriDecode = _evancz$elm_http$Native_Http.uriDecode;
var _evancz$elm_http$Http$uriEncode = _evancz$elm_http$Native_Http.uriEncode;
var _evancz$elm_http$Http$queryEscape = function (string) {
	return A2(
		_elm_lang$core$String$join,
		'+',
		A2(
			_elm_lang$core$String$split,
			'%20',
			_evancz$elm_http$Http$uriEncode(string)));
};
var _evancz$elm_http$Http$queryPair = function (_p0) {
	var _p1 = _p0;
	return A2(
		_elm_lang$core$Basics_ops['++'],
		_evancz$elm_http$Http$queryEscape(_p1._0),
		A2(
			_elm_lang$core$Basics_ops['++'],
			'=',
			_evancz$elm_http$Http$queryEscape(_p1._1)));
};
var _evancz$elm_http$Http$url = F2(
	function (baseUrl, args) {
		var _p2 = args;
		if (_p2.ctor === '[]') {
			return baseUrl;
		} else {
			return A2(
				_elm_lang$core$Basics_ops['++'],
				baseUrl,
				A2(
					_elm_lang$core$Basics_ops['++'],
					'?',
					A2(
						_elm_lang$core$String$join,
						'&',
						A2(_elm_lang$core$List$map, _evancz$elm_http$Http$queryPair, args))));
		}
	});
var _evancz$elm_http$Http$Request = F4(
	function (a, b, c, d) {
		return {verb: a, headers: b, url: c, body: d};
	});
var _evancz$elm_http$Http$Settings = F5(
	function (a, b, c, d, e) {
		return {timeout: a, onStart: b, onProgress: c, desiredResponseType: d, withCredentials: e};
	});
var _evancz$elm_http$Http$Response = F5(
	function (a, b, c, d, e) {
		return {status: a, statusText: b, headers: c, url: d, value: e};
	});
var _evancz$elm_http$Http$TODO_implement_blob_in_another_library = {ctor: 'TODO_implement_blob_in_another_library'};
var _evancz$elm_http$Http$TODO_implement_file_in_another_library = {ctor: 'TODO_implement_file_in_another_library'};
var _evancz$elm_http$Http$BodyBlob = function (a) {
	return {ctor: 'BodyBlob', _0: a};
};
var _evancz$elm_http$Http$BodyFormData = {ctor: 'BodyFormData'};
var _evancz$elm_http$Http$ArrayBuffer = {ctor: 'ArrayBuffer'};
var _evancz$elm_http$Http$BodyString = function (a) {
	return {ctor: 'BodyString', _0: a};
};
var _evancz$elm_http$Http$string = _evancz$elm_http$Http$BodyString;
var _evancz$elm_http$Http$Empty = {ctor: 'Empty'};
var _evancz$elm_http$Http$empty = _evancz$elm_http$Http$Empty;
var _evancz$elm_http$Http$FileData = F3(
	function (a, b, c) {
		return {ctor: 'FileData', _0: a, _1: b, _2: c};
	});
var _evancz$elm_http$Http$BlobData = F3(
	function (a, b, c) {
		return {ctor: 'BlobData', _0: a, _1: b, _2: c};
	});
var _evancz$elm_http$Http$blobData = _evancz$elm_http$Http$BlobData;
var _evancz$elm_http$Http$StringData = F2(
	function (a, b) {
		return {ctor: 'StringData', _0: a, _1: b};
	});
var _evancz$elm_http$Http$stringData = _evancz$elm_http$Http$StringData;
var _evancz$elm_http$Http$Blob = function (a) {
	return {ctor: 'Blob', _0: a};
};
var _evancz$elm_http$Http$Text = function (a) {
	return {ctor: 'Text', _0: a};
};
var _evancz$elm_http$Http$RawNetworkError = {ctor: 'RawNetworkError'};
var _evancz$elm_http$Http$RawTimeout = {ctor: 'RawTimeout'};
var _evancz$elm_http$Http$BadResponse = F2(
	function (a, b) {
		return {ctor: 'BadResponse', _0: a, _1: b};
	});
var _evancz$elm_http$Http$UnexpectedPayload = function (a) {
	return {ctor: 'UnexpectedPayload', _0: a};
};
var _evancz$elm_http$Http$handleResponse = F2(
	function (handle, response) {
		if ((_elm_lang$core$Native_Utils.cmp(200, response.status) < 1) && (_elm_lang$core$Native_Utils.cmp(response.status, 300) < 0)) {
			var _p3 = response.value;
			if (_p3.ctor === 'Text') {
				return handle(_p3._0);
			} else {
				return _elm_lang$core$Task$fail(
					_evancz$elm_http$Http$UnexpectedPayload('Response body is a blob, expecting a string.'));
			}
		} else {
			return _elm_lang$core$Task$fail(
				A2(_evancz$elm_http$Http$BadResponse, response.status, response.statusText));
		}
	});
var _evancz$elm_http$Http$NetworkError = {ctor: 'NetworkError'};
var _evancz$elm_http$Http$Timeout = {ctor: 'Timeout'};
var _evancz$elm_http$Http$promoteError = function (rawError) {
	var _p4 = rawError;
	if (_p4.ctor === 'RawTimeout') {
		return _evancz$elm_http$Http$Timeout;
	} else {
		return _evancz$elm_http$Http$NetworkError;
	}
};
var _evancz$elm_http$Http$getString = function (url) {
	var request = {
		verb: 'GET',
		headers: _elm_lang$core$Native_List.fromArray(
			[]),
		url: url,
		body: _evancz$elm_http$Http$empty
	};
	return A2(
		_elm_lang$core$Task$andThen,
		A2(
			_elm_lang$core$Task$mapError,
			_evancz$elm_http$Http$promoteError,
			A2(_evancz$elm_http$Http$send, _evancz$elm_http$Http$defaultSettings, request)),
		_evancz$elm_http$Http$handleResponse(_elm_lang$core$Task$succeed));
};
var _evancz$elm_http$Http$fromJson = F2(
	function (decoder, response) {
		var decode = function (str) {
			var _p5 = A2(_elm_lang$core$Json_Decode$decodeString, decoder, str);
			if (_p5.ctor === 'Ok') {
				return _elm_lang$core$Task$succeed(_p5._0);
			} else {
				return _elm_lang$core$Task$fail(
					_evancz$elm_http$Http$UnexpectedPayload(_p5._0));
			}
		};
		return A2(
			_elm_lang$core$Task$andThen,
			A2(_elm_lang$core$Task$mapError, _evancz$elm_http$Http$promoteError, response),
			_evancz$elm_http$Http$handleResponse(decode));
	});
var _evancz$elm_http$Http$get = F2(
	function (decoder, url) {
		var request = {
			verb: 'GET',
			headers: _elm_lang$core$Native_List.fromArray(
				[]),
			url: url,
			body: _evancz$elm_http$Http$empty
		};
		return A2(
			_evancz$elm_http$Http$fromJson,
			decoder,
			A2(_evancz$elm_http$Http$send, _evancz$elm_http$Http$defaultSettings, request));
	});
var _evancz$elm_http$Http$post = F3(
	function (decoder, url, body) {
		var request = {
			verb: 'POST',
			headers: _elm_lang$core$Native_List.fromArray(
				[]),
			url: url,
			body: body
		};
		return A2(
			_evancz$elm_http$Http$fromJson,
			decoder,
			A2(_evancz$elm_http$Http$send, _evancz$elm_http$Http$defaultSettings, request));
	});

var _elm_lang$navigation$Native_Navigation = function() {

function go(n)
{
	return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback)
	{
		if (n !== 0)
		{
			history.go(n);
		}
		callback(_elm_lang$core$Native_Scheduler.succeed(_elm_lang$core$Native_Utils.Tuple0));
	});
}

function pushState(url)
{
	return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback)
	{
		history.pushState({}, '', url);
		callback(_elm_lang$core$Native_Scheduler.succeed(getLocation()));
	});
}

function replaceState(url)
{
	return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback)
	{
		history.replaceState({}, '', url);
		callback(_elm_lang$core$Native_Scheduler.succeed(getLocation()));
	});
}

function getLocation()
{
	var location = document.location;

	return {
		href: location.href,
		host: location.host,
		hostname: location.hostname,
		protocol: location.protocol,
		origin: location.origin,
		port_: location.port,
		pathname: location.pathname,
		search: location.search,
		hash: location.hash,
		username: location.username,
		password: location.password
	};
}


return {
	go: go,
	pushState: pushState,
	replaceState: replaceState,
	getLocation: getLocation
};

}();

var _elm_lang$navigation$Navigation$replaceState = _elm_lang$navigation$Native_Navigation.replaceState;
var _elm_lang$navigation$Navigation$pushState = _elm_lang$navigation$Native_Navigation.pushState;
var _elm_lang$navigation$Navigation$go = _elm_lang$navigation$Native_Navigation.go;
var _elm_lang$navigation$Navigation$spawnPopState = function (router) {
	return _elm_lang$core$Process$spawn(
		A3(
			_elm_lang$dom$Dom_LowLevel$onWindow,
			'popstate',
			_elm_lang$core$Json_Decode$value,
			function (_p0) {
				return A2(
					_elm_lang$core$Platform$sendToSelf,
					router,
					_elm_lang$navigation$Native_Navigation.getLocation(
						{ctor: '_Tuple0'}));
			}));
};
var _elm_lang$navigation$Navigation_ops = _elm_lang$navigation$Navigation_ops || {};
_elm_lang$navigation$Navigation_ops['&>'] = F2(
	function (task1, task2) {
		return A2(
			_elm_lang$core$Task$andThen,
			task1,
			function (_p1) {
				return task2;
			});
	});
var _elm_lang$navigation$Navigation$notify = F3(
	function (router, subs, location) {
		var send = function (_p2) {
			var _p3 = _p2;
			return A2(
				_elm_lang$core$Platform$sendToApp,
				router,
				_p3._0(location));
		};
		return A2(
			_elm_lang$navigation$Navigation_ops['&>'],
			_elm_lang$core$Task$sequence(
				A2(_elm_lang$core$List$map, send, subs)),
			_elm_lang$core$Task$succeed(
				{ctor: '_Tuple0'}));
	});
var _elm_lang$navigation$Navigation$onSelfMsg = F3(
	function (router, location, state) {
		return A2(
			_elm_lang$navigation$Navigation_ops['&>'],
			A3(_elm_lang$navigation$Navigation$notify, router, state.subs, location),
			_elm_lang$core$Task$succeed(state));
	});
var _elm_lang$navigation$Navigation$cmdHelp = F3(
	function (router, subs, cmd) {
		var _p4 = cmd;
		switch (_p4.ctor) {
			case 'Jump':
				return _elm_lang$navigation$Navigation$go(_p4._0);
			case 'New':
				return A2(
					_elm_lang$core$Task$andThen,
					_elm_lang$navigation$Navigation$pushState(_p4._0),
					A2(_elm_lang$navigation$Navigation$notify, router, subs));
			default:
				return A2(
					_elm_lang$core$Task$andThen,
					_elm_lang$navigation$Navigation$replaceState(_p4._0),
					A2(_elm_lang$navigation$Navigation$notify, router, subs));
		}
	});
var _elm_lang$navigation$Navigation$updateHelp = F2(
	function (func, _p5) {
		var _p6 = _p5;
		return {
			ctor: '_Tuple2',
			_0: _p6._0,
			_1: A2(_elm_lang$core$Platform_Cmd$map, func, _p6._1)
		};
	});
var _elm_lang$navigation$Navigation$subscription = _elm_lang$core$Native_Platform.leaf('Navigation');
var _elm_lang$navigation$Navigation$command = _elm_lang$core$Native_Platform.leaf('Navigation');
var _elm_lang$navigation$Navigation$Location = function (a) {
	return function (b) {
		return function (c) {
			return function (d) {
				return function (e) {
					return function (f) {
						return function (g) {
							return function (h) {
								return function (i) {
									return function (j) {
										return function (k) {
											return {href: a, host: b, hostname: c, protocol: d, origin: e, port_: f, pathname: g, search: h, hash: i, username: j, password: k};
										};
									};
								};
							};
						};
					};
				};
			};
		};
	};
};
var _elm_lang$navigation$Navigation$State = F2(
	function (a, b) {
		return {subs: a, process: b};
	});
var _elm_lang$navigation$Navigation$init = _elm_lang$core$Task$succeed(
	A2(
		_elm_lang$navigation$Navigation$State,
		_elm_lang$core$Native_List.fromArray(
			[]),
		_elm_lang$core$Maybe$Nothing));
var _elm_lang$navigation$Navigation$onEffects = F4(
	function (router, cmds, subs, _p7) {
		var _p8 = _p7;
		var _p10 = _p8.process;
		var stepState = function () {
			var _p9 = {ctor: '_Tuple2', _0: subs, _1: _p10};
			_v4_2:
			do {
				if (_p9._0.ctor === '[]') {
					if (_p9._1.ctor === 'Just') {
						return A2(
							_elm_lang$navigation$Navigation_ops['&>'],
							_elm_lang$core$Process$kill(_p9._1._0),
							_elm_lang$core$Task$succeed(
								A2(_elm_lang$navigation$Navigation$State, subs, _elm_lang$core$Maybe$Nothing)));
					} else {
						break _v4_2;
					}
				} else {
					if (_p9._1.ctor === 'Nothing') {
						return A2(
							_elm_lang$core$Task$andThen,
							_elm_lang$navigation$Navigation$spawnPopState(router),
							function (pid) {
								return _elm_lang$core$Task$succeed(
									A2(
										_elm_lang$navigation$Navigation$State,
										subs,
										_elm_lang$core$Maybe$Just(pid)));
							});
					} else {
						break _v4_2;
					}
				}
			} while(false);
			return _elm_lang$core$Task$succeed(
				A2(_elm_lang$navigation$Navigation$State, subs, _p10));
		}();
		return A2(
			_elm_lang$navigation$Navigation_ops['&>'],
			_elm_lang$core$Task$sequence(
				A2(
					_elm_lang$core$List$map,
					A2(_elm_lang$navigation$Navigation$cmdHelp, router, subs),
					cmds)),
			stepState);
	});
var _elm_lang$navigation$Navigation$UserMsg = function (a) {
	return {ctor: 'UserMsg', _0: a};
};
var _elm_lang$navigation$Navigation$Change = function (a) {
	return {ctor: 'Change', _0: a};
};
var _elm_lang$navigation$Navigation$Parser = function (a) {
	return {ctor: 'Parser', _0: a};
};
var _elm_lang$navigation$Navigation$makeParser = _elm_lang$navigation$Navigation$Parser;
var _elm_lang$navigation$Navigation$Modify = function (a) {
	return {ctor: 'Modify', _0: a};
};
var _elm_lang$navigation$Navigation$modifyUrl = function (url) {
	return _elm_lang$navigation$Navigation$command(
		_elm_lang$navigation$Navigation$Modify(url));
};
var _elm_lang$navigation$Navigation$New = function (a) {
	return {ctor: 'New', _0: a};
};
var _elm_lang$navigation$Navigation$newUrl = function (url) {
	return _elm_lang$navigation$Navigation$command(
		_elm_lang$navigation$Navigation$New(url));
};
var _elm_lang$navigation$Navigation$Jump = function (a) {
	return {ctor: 'Jump', _0: a};
};
var _elm_lang$navigation$Navigation$back = function (n) {
	return _elm_lang$navigation$Navigation$command(
		_elm_lang$navigation$Navigation$Jump(0 - n));
};
var _elm_lang$navigation$Navigation$forward = function (n) {
	return _elm_lang$navigation$Navigation$command(
		_elm_lang$navigation$Navigation$Jump(n));
};
var _elm_lang$navigation$Navigation$cmdMap = F2(
	function (_p11, myCmd) {
		var _p12 = myCmd;
		switch (_p12.ctor) {
			case 'Jump':
				return _elm_lang$navigation$Navigation$Jump(_p12._0);
			case 'New':
				return _elm_lang$navigation$Navigation$New(_p12._0);
			default:
				return _elm_lang$navigation$Navigation$Modify(_p12._0);
		}
	});
var _elm_lang$navigation$Navigation$Monitor = function (a) {
	return {ctor: 'Monitor', _0: a};
};
var _elm_lang$navigation$Navigation$programWithFlags = F2(
	function (_p13, stuff) {
		var _p14 = _p13;
		var _p16 = _p14._0;
		var location = _elm_lang$navigation$Native_Navigation.getLocation(
			{ctor: '_Tuple0'});
		var init = function (flags) {
			return A2(
				_elm_lang$navigation$Navigation$updateHelp,
				_elm_lang$navigation$Navigation$UserMsg,
				A2(
					stuff.init,
					flags,
					_p16(location)));
		};
		var view = function (model) {
			return A2(
				_elm_lang$html$Html_App$map,
				_elm_lang$navigation$Navigation$UserMsg,
				stuff.view(model));
		};
		var subs = function (model) {
			return _elm_lang$core$Platform_Sub$batch(
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$navigation$Navigation$subscription(
						_elm_lang$navigation$Navigation$Monitor(_elm_lang$navigation$Navigation$Change)),
						A2(
						_elm_lang$core$Platform_Sub$map,
						_elm_lang$navigation$Navigation$UserMsg,
						stuff.subscriptions(model))
					]));
		};
		var update = F2(
			function (msg, model) {
				return A2(
					_elm_lang$navigation$Navigation$updateHelp,
					_elm_lang$navigation$Navigation$UserMsg,
					function () {
						var _p15 = msg;
						if (_p15.ctor === 'Change') {
							return A2(
								stuff.urlUpdate,
								_p16(_p15._0),
								model);
						} else {
							return A2(stuff.update, _p15._0, model);
						}
					}());
			});
		return _elm_lang$html$Html_App$programWithFlags(
			{init: init, view: view, update: update, subscriptions: subs});
	});
var _elm_lang$navigation$Navigation$program = F2(
	function (parser, stuff) {
		return A2(
			_elm_lang$navigation$Navigation$programWithFlags,
			parser,
			_elm_lang$core$Native_Utils.update(
				stuff,
				{
					init: function (_p17) {
						return stuff.init;
					}
				}));
	});
var _elm_lang$navigation$Navigation$subMap = F2(
	function (func, _p18) {
		var _p19 = _p18;
		return _elm_lang$navigation$Navigation$Monitor(
			function (_p20) {
				return func(
					_p19._0(_p20));
			});
	});
_elm_lang$core$Native_Platform.effectManagers['Navigation'] = {pkg: 'elm-lang/navigation', init: _elm_lang$navigation$Navigation$init, onEffects: _elm_lang$navigation$Navigation$onEffects, onSelfMsg: _elm_lang$navigation$Navigation$onSelfMsg, tag: 'fx', cmdMap: _elm_lang$navigation$Navigation$cmdMap, subMap: _elm_lang$navigation$Navigation$subMap};

var _vito$elm_ansi$Ansi$encodeCode = function (code) {
	var _p0 = code;
	if (_p0.ctor === 'Nothing') {
		return '';
	} else {
		return _elm_lang$core$Basics$toString(_p0._0);
	}
};
var _vito$elm_ansi$Ansi$encodeCodes = function (codes) {
	return A2(
		_elm_lang$core$String$join,
		';',
		A2(_elm_lang$core$List$map, _vito$elm_ansi$Ansi$encodeCode, codes));
};
var _vito$elm_ansi$Ansi$RestoreCursorPosition = {ctor: 'RestoreCursorPosition'};
var _vito$elm_ansi$Ansi$SaveCursorPosition = {ctor: 'SaveCursorPosition'};
var _vito$elm_ansi$Ansi$EraseLine = function (a) {
	return {ctor: 'EraseLine', _0: a};
};
var _vito$elm_ansi$Ansi$EraseDisplay = function (a) {
	return {ctor: 'EraseDisplay', _0: a};
};
var _vito$elm_ansi$Ansi$CursorColumn = function (a) {
	return {ctor: 'CursorColumn', _0: a};
};
var _vito$elm_ansi$Ansi$CursorPosition = F2(
	function (a, b) {
		return {ctor: 'CursorPosition', _0: a, _1: b};
	});
var _vito$elm_ansi$Ansi$cursorPosition = function (codes) {
	var _p1 = codes;
	_v1_5:
	do {
		if (_p1.ctor === '::') {
			if (_p1._0.ctor === 'Nothing') {
				if (_p1._1.ctor === '[]') {
					return _elm_lang$core$Native_List.fromArray(
						[
							A2(_vito$elm_ansi$Ansi$CursorPosition, 1, 1)
						]);
				} else {
					if (_p1._1._0.ctor === 'Nothing') {
						if (_p1._1._1.ctor === '[]') {
							return _elm_lang$core$Native_List.fromArray(
								[
									A2(_vito$elm_ansi$Ansi$CursorPosition, 1, 1)
								]);
						} else {
							break _v1_5;
						}
					} else {
						if (_p1._1._1.ctor === '[]') {
							return _elm_lang$core$Native_List.fromArray(
								[
									A2(_vito$elm_ansi$Ansi$CursorPosition, 1, _p1._1._0._0)
								]);
						} else {
							break _v1_5;
						}
					}
				}
			} else {
				if (_p1._1.ctor === '::') {
					if (_p1._1._0.ctor === 'Nothing') {
						if (_p1._1._1.ctor === '[]') {
							return _elm_lang$core$Native_List.fromArray(
								[
									A2(_vito$elm_ansi$Ansi$CursorPosition, _p1._0._0, 1)
								]);
						} else {
							break _v1_5;
						}
					} else {
						if (_p1._1._1.ctor === '[]') {
							return _elm_lang$core$Native_List.fromArray(
								[
									A2(_vito$elm_ansi$Ansi$CursorPosition, _p1._0._0, _p1._1._0._0)
								]);
						} else {
							break _v1_5;
						}
					}
				} else {
					break _v1_5;
				}
			}
		} else {
			break _v1_5;
		}
	} while(false);
	return _elm_lang$core$Native_List.fromArray(
		[]);
};
var _vito$elm_ansi$Ansi$CursorBack = function (a) {
	return {ctor: 'CursorBack', _0: a};
};
var _vito$elm_ansi$Ansi$CursorForward = function (a) {
	return {ctor: 'CursorForward', _0: a};
};
var _vito$elm_ansi$Ansi$CursorDown = function (a) {
	return {ctor: 'CursorDown', _0: a};
};
var _vito$elm_ansi$Ansi$CursorUp = function (a) {
	return {ctor: 'CursorUp', _0: a};
};
var _vito$elm_ansi$Ansi$CarriageReturn = {ctor: 'CarriageReturn'};
var _vito$elm_ansi$Ansi$Linebreak = {ctor: 'Linebreak'};
var _vito$elm_ansi$Ansi$SetInverted = function (a) {
	return {ctor: 'SetInverted', _0: a};
};
var _vito$elm_ansi$Ansi$SetBlink = function (a) {
	return {ctor: 'SetBlink', _0: a};
};
var _vito$elm_ansi$Ansi$SetUnderline = function (a) {
	return {ctor: 'SetUnderline', _0: a};
};
var _vito$elm_ansi$Ansi$SetItalic = function (a) {
	return {ctor: 'SetItalic', _0: a};
};
var _vito$elm_ansi$Ansi$SetFaint = function (a) {
	return {ctor: 'SetFaint', _0: a};
};
var _vito$elm_ansi$Ansi$SetBold = function (a) {
	return {ctor: 'SetBold', _0: a};
};
var _vito$elm_ansi$Ansi$SetBackground = function (a) {
	return {ctor: 'SetBackground', _0: a};
};
var _vito$elm_ansi$Ansi$SetForeground = function (a) {
	return {ctor: 'SetForeground', _0: a};
};
var _vito$elm_ansi$Ansi$reset = _elm_lang$core$Native_List.fromArray(
	[
		_vito$elm_ansi$Ansi$SetForeground(_elm_lang$core$Maybe$Nothing),
		_vito$elm_ansi$Ansi$SetBackground(_elm_lang$core$Maybe$Nothing),
		_vito$elm_ansi$Ansi$SetBold(false),
		_vito$elm_ansi$Ansi$SetFaint(false),
		_vito$elm_ansi$Ansi$SetItalic(false),
		_vito$elm_ansi$Ansi$SetUnderline(false),
		_vito$elm_ansi$Ansi$SetBlink(false),
		_vito$elm_ansi$Ansi$SetInverted(false)
	]);
var _vito$elm_ansi$Ansi$Remainder = function (a) {
	return {ctor: 'Remainder', _0: a};
};
var _vito$elm_ansi$Ansi$Print = function (a) {
	return {ctor: 'Print', _0: a};
};
var _vito$elm_ansi$Ansi$completeParsing = function (parser) {
	var _p2 = parser;
	switch (_p2._0.ctor) {
		case 'Escaped':
			return A2(
				_p2._2,
				_vito$elm_ansi$Ansi$Remainder(''),
				_p2._1);
		case 'CSI':
			return A2(
				_p2._2,
				_vito$elm_ansi$Ansi$Remainder(
					A2(
						_elm_lang$core$Basics_ops['++'],
						'[',
						_vito$elm_ansi$Ansi$encodeCodes(
							A2(
								_elm_lang$core$Basics_ops['++'],
								_p2._0._0,
								_elm_lang$core$Native_List.fromArray(
									[_p2._0._1]))))),
				_p2._1);
		default:
			if (_p2._0._0 === '') {
				return _p2._1;
			} else {
				return A2(
					_p2._2,
					_vito$elm_ansi$Ansi$Print(_p2._0._0),
					_p2._1);
			}
	}
};
var _vito$elm_ansi$Ansi$completeUnescaped = function (parser) {
	var _p3 = parser;
	if (_p3._0.ctor === 'Unescaped') {
		if (_p3._0._0 === '') {
			return _p3._1;
		} else {
			return A2(
				_p3._2,
				_vito$elm_ansi$Ansi$Print(_p3._0._0),
				_p3._1);
		}
	} else {
		return _p3._1;
	}
};
var _vito$elm_ansi$Ansi$BrightWhite = {ctor: 'BrightWhite'};
var _vito$elm_ansi$Ansi$BrightCyan = {ctor: 'BrightCyan'};
var _vito$elm_ansi$Ansi$BrightMagenta = {ctor: 'BrightMagenta'};
var _vito$elm_ansi$Ansi$BrightBlue = {ctor: 'BrightBlue'};
var _vito$elm_ansi$Ansi$BrightYellow = {ctor: 'BrightYellow'};
var _vito$elm_ansi$Ansi$BrightGreen = {ctor: 'BrightGreen'};
var _vito$elm_ansi$Ansi$BrightRed = {ctor: 'BrightRed'};
var _vito$elm_ansi$Ansi$BrightBlack = {ctor: 'BrightBlack'};
var _vito$elm_ansi$Ansi$White = {ctor: 'White'};
var _vito$elm_ansi$Ansi$Cyan = {ctor: 'Cyan'};
var _vito$elm_ansi$Ansi$Magenta = {ctor: 'Magenta'};
var _vito$elm_ansi$Ansi$Blue = {ctor: 'Blue'};
var _vito$elm_ansi$Ansi$Yellow = {ctor: 'Yellow'};
var _vito$elm_ansi$Ansi$Green = {ctor: 'Green'};
var _vito$elm_ansi$Ansi$Red = {ctor: 'Red'};
var _vito$elm_ansi$Ansi$Black = {ctor: 'Black'};
var _vito$elm_ansi$Ansi$codeActions = function (code) {
	var _p4 = code;
	switch (_p4) {
		case 0:
			return _vito$elm_ansi$Ansi$reset;
		case 1:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBold(true)
				]);
		case 2:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetFaint(true)
				]);
		case 3:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetItalic(true)
				]);
		case 4:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetUnderline(true)
				]);
		case 5:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBlink(true)
				]);
		case 7:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetInverted(true)
				]);
		case 30:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetForeground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$Black))
				]);
		case 31:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetForeground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$Red))
				]);
		case 32:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetForeground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$Green))
				]);
		case 33:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetForeground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$Yellow))
				]);
		case 34:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetForeground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$Blue))
				]);
		case 35:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetForeground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$Magenta))
				]);
		case 36:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetForeground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$Cyan))
				]);
		case 37:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetForeground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$White))
				]);
		case 39:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetForeground(_elm_lang$core$Maybe$Nothing)
				]);
		case 40:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBackground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$Black))
				]);
		case 41:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBackground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$Red))
				]);
		case 42:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBackground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$Green))
				]);
		case 43:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBackground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$Yellow))
				]);
		case 44:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBackground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$Blue))
				]);
		case 45:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBackground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$Magenta))
				]);
		case 46:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBackground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$Cyan))
				]);
		case 47:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBackground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$White))
				]);
		case 49:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBackground(_elm_lang$core$Maybe$Nothing)
				]);
		case 90:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetForeground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$BrightBlack))
				]);
		case 91:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetForeground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$BrightRed))
				]);
		case 92:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetForeground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$BrightGreen))
				]);
		case 93:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetForeground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$BrightYellow))
				]);
		case 94:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetForeground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$BrightBlue))
				]);
		case 95:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetForeground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$BrightMagenta))
				]);
		case 96:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetForeground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$BrightCyan))
				]);
		case 97:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetForeground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$BrightWhite))
				]);
		case 100:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBackground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$BrightBlack))
				]);
		case 101:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBackground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$BrightRed))
				]);
		case 102:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBackground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$BrightGreen))
				]);
		case 103:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBackground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$BrightYellow))
				]);
		case 104:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBackground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$BrightBlue))
				]);
		case 105:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBackground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$BrightMagenta))
				]);
		case 106:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBackground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$BrightCyan))
				]);
		case 107:
			return _elm_lang$core$Native_List.fromArray(
				[
					_vito$elm_ansi$Ansi$SetBackground(
					_elm_lang$core$Maybe$Just(_vito$elm_ansi$Ansi$BrightWhite))
				]);
		default:
			return _elm_lang$core$Native_List.fromArray(
				[]);
	}
};
var _vito$elm_ansi$Ansi$EraseAll = {ctor: 'EraseAll'};
var _vito$elm_ansi$Ansi$EraseToEnd = {ctor: 'EraseToEnd'};
var _vito$elm_ansi$Ansi$EraseToBeginning = {ctor: 'EraseToBeginning'};
var _vito$elm_ansi$Ansi$eraseMode = function (code) {
	var _p5 = code;
	switch (_p5) {
		case 0:
			return _vito$elm_ansi$Ansi$EraseToEnd;
		case 1:
			return _vito$elm_ansi$Ansi$EraseToBeginning;
		default:
			return _vito$elm_ansi$Ansi$EraseAll;
	}
};
var _vito$elm_ansi$Ansi$Parser = F3(
	function (a, b, c) {
		return {ctor: 'Parser', _0: a, _1: b, _2: c};
	});
var _vito$elm_ansi$Ansi$Unescaped = function (a) {
	return {ctor: 'Unescaped', _0: a};
};
var _vito$elm_ansi$Ansi$emptyParser = _vito$elm_ansi$Ansi$Parser(
	_vito$elm_ansi$Ansi$Unescaped(''));
var _vito$elm_ansi$Ansi$completeBracketed = F2(
	function (_p6, actions) {
		var _p7 = _p6;
		var _p8 = _p7._2;
		return A3(
			_vito$elm_ansi$Ansi$Parser,
			_vito$elm_ansi$Ansi$Unescaped(''),
			A3(_elm_lang$core$List$foldl, _p8, _p7._1, actions),
			_p8);
	});
var _vito$elm_ansi$Ansi$CSI = F2(
	function (a, b) {
		return {ctor: 'CSI', _0: a, _1: b};
	});
var _vito$elm_ansi$Ansi$Escaped = {ctor: 'Escaped'};
var _vito$elm_ansi$Ansi$parseChar = F2(
	function ($char, parser) {
		var _p9 = parser;
		switch (_p9._0.ctor) {
			case 'Unescaped':
				var _p11 = _p9._2;
				var _p10 = $char;
				switch (_p10) {
					case '\r':
						return A3(
							_vito$elm_ansi$Ansi$Parser,
							_vito$elm_ansi$Ansi$Unescaped(''),
							A2(
								_p11,
								_vito$elm_ansi$Ansi$CarriageReturn,
								_vito$elm_ansi$Ansi$completeUnescaped(parser)),
							_p11);
					case '\n':
						return A3(
							_vito$elm_ansi$Ansi$Parser,
							_vito$elm_ansi$Ansi$Unescaped(''),
							A2(
								_p11,
								_vito$elm_ansi$Ansi$Linebreak,
								_vito$elm_ansi$Ansi$completeUnescaped(parser)),
							_p11);
					case '':
						return A3(
							_vito$elm_ansi$Ansi$Parser,
							_vito$elm_ansi$Ansi$Escaped,
							_vito$elm_ansi$Ansi$completeUnescaped(parser),
							_p11);
					default:
						return A3(
							_vito$elm_ansi$Ansi$Parser,
							_vito$elm_ansi$Ansi$Unescaped(
								A2(_elm_lang$core$Basics_ops['++'], _p9._0._0, $char)),
							_p9._1,
							_p11);
				}
			case 'Escaped':
				var _p14 = _p9._2;
				var _p13 = _p9._1;
				var _p12 = $char;
				if (_p12 === '[') {
					return A3(
						_vito$elm_ansi$Ansi$Parser,
						A2(
							_vito$elm_ansi$Ansi$CSI,
							_elm_lang$core$Native_List.fromArray(
								[]),
							_elm_lang$core$Maybe$Nothing),
						_p13,
						_p14);
				} else {
					return A3(
						_vito$elm_ansi$Ansi$Parser,
						_vito$elm_ansi$Ansi$Unescaped($char),
						_p13,
						_p14);
				}
			default:
				var _p21 = _p9._2;
				var _p20 = _p9._1;
				var _p19 = _p9._0._1;
				var _p18 = _p9._0._0;
				var _p15 = $char;
				switch (_p15) {
					case 'm':
						return A2(
							_vito$elm_ansi$Ansi$completeBracketed,
							parser,
							A2(
								_elm_lang$core$List$concatMap,
								function (_p16) {
									return _vito$elm_ansi$Ansi$codeActions(
										A2(_elm_lang$core$Maybe$withDefault, 0, _p16));
								},
								A2(
									_elm_lang$core$Basics_ops['++'],
									_p18,
									_elm_lang$core$Native_List.fromArray(
										[_p19]))));
					case 'A':
						return A2(
							_vito$elm_ansi$Ansi$completeBracketed,
							parser,
							_elm_lang$core$Native_List.fromArray(
								[
									_vito$elm_ansi$Ansi$CursorUp(
									A2(_elm_lang$core$Maybe$withDefault, 1, _p19))
								]));
					case 'B':
						return A2(
							_vito$elm_ansi$Ansi$completeBracketed,
							parser,
							_elm_lang$core$Native_List.fromArray(
								[
									_vito$elm_ansi$Ansi$CursorDown(
									A2(_elm_lang$core$Maybe$withDefault, 1, _p19))
								]));
					case 'C':
						return A2(
							_vito$elm_ansi$Ansi$completeBracketed,
							parser,
							_elm_lang$core$Native_List.fromArray(
								[
									_vito$elm_ansi$Ansi$CursorForward(
									A2(_elm_lang$core$Maybe$withDefault, 1, _p19))
								]));
					case 'D':
						return A2(
							_vito$elm_ansi$Ansi$completeBracketed,
							parser,
							_elm_lang$core$Native_List.fromArray(
								[
									_vito$elm_ansi$Ansi$CursorBack(
									A2(_elm_lang$core$Maybe$withDefault, 1, _p19))
								]));
					case 'E':
						return A2(
							_vito$elm_ansi$Ansi$completeBracketed,
							parser,
							_elm_lang$core$Native_List.fromArray(
								[
									_vito$elm_ansi$Ansi$CursorDown(
									A2(_elm_lang$core$Maybe$withDefault, 1, _p19)),
									_vito$elm_ansi$Ansi$CursorColumn(0)
								]));
					case 'F':
						return A2(
							_vito$elm_ansi$Ansi$completeBracketed,
							parser,
							_elm_lang$core$Native_List.fromArray(
								[
									_vito$elm_ansi$Ansi$CursorUp(
									A2(_elm_lang$core$Maybe$withDefault, 1, _p19)),
									_vito$elm_ansi$Ansi$CursorColumn(0)
								]));
					case 'G':
						return A2(
							_vito$elm_ansi$Ansi$completeBracketed,
							parser,
							_elm_lang$core$Native_List.fromArray(
								[
									_vito$elm_ansi$Ansi$CursorColumn(
									A2(_elm_lang$core$Maybe$withDefault, 0, _p19))
								]));
					case 'H':
						return A2(
							_vito$elm_ansi$Ansi$completeBracketed,
							parser,
							_vito$elm_ansi$Ansi$cursorPosition(
								A2(
									_elm_lang$core$Basics_ops['++'],
									_p18,
									_elm_lang$core$Native_List.fromArray(
										[_p19]))));
					case 'J':
						return A2(
							_vito$elm_ansi$Ansi$completeBracketed,
							parser,
							_elm_lang$core$Native_List.fromArray(
								[
									_vito$elm_ansi$Ansi$EraseDisplay(
									_vito$elm_ansi$Ansi$eraseMode(
										A2(_elm_lang$core$Maybe$withDefault, 0, _p19)))
								]));
					case 'K':
						return A2(
							_vito$elm_ansi$Ansi$completeBracketed,
							parser,
							_elm_lang$core$Native_List.fromArray(
								[
									_vito$elm_ansi$Ansi$EraseLine(
									_vito$elm_ansi$Ansi$eraseMode(
										A2(_elm_lang$core$Maybe$withDefault, 0, _p19)))
								]));
					case 'f':
						return A2(
							_vito$elm_ansi$Ansi$completeBracketed,
							parser,
							_vito$elm_ansi$Ansi$cursorPosition(
								A2(
									_elm_lang$core$Basics_ops['++'],
									_p18,
									_elm_lang$core$Native_List.fromArray(
										[_p19]))));
					case 's':
						return A2(
							_vito$elm_ansi$Ansi$completeBracketed,
							parser,
							_elm_lang$core$Native_List.fromArray(
								[_vito$elm_ansi$Ansi$SaveCursorPosition]));
					case 'u':
						return A2(
							_vito$elm_ansi$Ansi$completeBracketed,
							parser,
							_elm_lang$core$Native_List.fromArray(
								[_vito$elm_ansi$Ansi$RestoreCursorPosition]));
					case ';':
						return A3(
							_vito$elm_ansi$Ansi$Parser,
							A2(
								_vito$elm_ansi$Ansi$CSI,
								A2(
									_elm_lang$core$Basics_ops['++'],
									_p18,
									_elm_lang$core$Native_List.fromArray(
										[_p19])),
								_elm_lang$core$Maybe$Nothing),
							_p20,
							_p21);
					default:
						var _p17 = _elm_lang$core$String$toInt(_p15);
						if (_p17.ctor === 'Ok') {
							return A3(
								_vito$elm_ansi$Ansi$Parser,
								A2(
									_vito$elm_ansi$Ansi$CSI,
									_p18,
									_elm_lang$core$Maybe$Just(
										(A2(_elm_lang$core$Maybe$withDefault, 0, _p19) * 10) + _p17._0)),
								_p20,
								_p21);
						} else {
							return A2(
								_vito$elm_ansi$Ansi$completeBracketed,
								parser,
								_elm_lang$core$Native_List.fromArray(
									[]));
						}
				}
		}
	});
var _vito$elm_ansi$Ansi$parseInto = F3(
	function (model, update, ansi) {
		return _vito$elm_ansi$Ansi$completeParsing(
			A3(
				_elm_lang$core$List$foldl,
				_vito$elm_ansi$Ansi$parseChar,
				A2(_vito$elm_ansi$Ansi$emptyParser, model, update),
				A2(_elm_lang$core$String$split, '', ansi)));
	});
var _vito$elm_ansi$Ansi$parse = function (_p22) {
	return _elm_lang$core$List$reverse(
		A3(
			_vito$elm_ansi$Ansi$parseInto,
			_elm_lang$core$Native_List.fromArray(
				[]),
			F2(
				function (x, y) {
					return A2(_elm_lang$core$List_ops['::'], x, y);
				}),
			_p22));
};

var _vito$elm_ansi$Ansi_Log$colorClasses = F3(
	function (suffix, bold, mc) {
		var brightPrefix = 'ansi-bright-';
		var prefix = bold ? brightPrefix : 'ansi-';
		var _p0 = mc;
		if (_p0.ctor === 'Nothing') {
			return bold ? _elm_lang$core$Native_List.fromArray(
				['ansi-bold']) : _elm_lang$core$Native_List.fromArray(
				[]);
		} else {
			switch (_p0._0.ctor) {
				case 'Black':
					return _elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$core$Basics_ops['++'],
							prefix,
							A2(_elm_lang$core$Basics_ops['++'], 'black', suffix))
						]);
				case 'Red':
					return _elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$core$Basics_ops['++'],
							prefix,
							A2(_elm_lang$core$Basics_ops['++'], 'red', suffix))
						]);
				case 'Green':
					return _elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$core$Basics_ops['++'],
							prefix,
							A2(_elm_lang$core$Basics_ops['++'], 'green', suffix))
						]);
				case 'Yellow':
					return _elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$core$Basics_ops['++'],
							prefix,
							A2(_elm_lang$core$Basics_ops['++'], 'yellow', suffix))
						]);
				case 'Blue':
					return _elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$core$Basics_ops['++'],
							prefix,
							A2(_elm_lang$core$Basics_ops['++'], 'blue', suffix))
						]);
				case 'Magenta':
					return _elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$core$Basics_ops['++'],
							prefix,
							A2(_elm_lang$core$Basics_ops['++'], 'magenta', suffix))
						]);
				case 'Cyan':
					return _elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$core$Basics_ops['++'],
							prefix,
							A2(_elm_lang$core$Basics_ops['++'], 'cyan', suffix))
						]);
				case 'White':
					return _elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$core$Basics_ops['++'],
							prefix,
							A2(_elm_lang$core$Basics_ops['++'], 'white', suffix))
						]);
				case 'BrightBlack':
					return _elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$core$Basics_ops['++'],
							brightPrefix,
							A2(_elm_lang$core$Basics_ops['++'], 'black', suffix))
						]);
				case 'BrightRed':
					return _elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$core$Basics_ops['++'],
							brightPrefix,
							A2(_elm_lang$core$Basics_ops['++'], 'red', suffix))
						]);
				case 'BrightGreen':
					return _elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$core$Basics_ops['++'],
							brightPrefix,
							A2(_elm_lang$core$Basics_ops['++'], 'green', suffix))
						]);
				case 'BrightYellow':
					return _elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$core$Basics_ops['++'],
							brightPrefix,
							A2(_elm_lang$core$Basics_ops['++'], 'yellow', suffix))
						]);
				case 'BrightBlue':
					return _elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$core$Basics_ops['++'],
							brightPrefix,
							A2(_elm_lang$core$Basics_ops['++'], 'blue', suffix))
						]);
				case 'BrightMagenta':
					return _elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$core$Basics_ops['++'],
							brightPrefix,
							A2(_elm_lang$core$Basics_ops['++'], 'magenta', suffix))
						]);
				case 'BrightCyan':
					return _elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$core$Basics_ops['++'],
							brightPrefix,
							A2(_elm_lang$core$Basics_ops['++'], 'cyan', suffix))
						]);
				default:
					return _elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$core$Basics_ops['++'],
							brightPrefix,
							A2(_elm_lang$core$Basics_ops['++'], 'white', suffix))
						]);
			}
		}
	});
var _vito$elm_ansi$Ansi_Log$styleAttributes = function (style) {
	return _elm_lang$core$Native_List.fromArray(
		[
			_elm_lang$html$Html_Attributes$style(
			_elm_lang$core$Native_List.fromArray(
				[
					{
					ctor: '_Tuple2',
					_0: 'font-weight',
					_1: style.bold ? 'bold' : 'normal'
				},
					{
					ctor: '_Tuple2',
					_0: 'text-decoration',
					_1: style.underline ? 'underline' : 'none'
				},
					{
					ctor: '_Tuple2',
					_0: 'font-style',
					_1: style.italic ? 'italic' : 'normal'
				}
				])),
			function () {
			var ansiClasses = _elm_lang$core$Native_List.fromArray(
				[
					{ctor: '_Tuple2', _0: 'ansi-blink', _1: style.blink},
					{ctor: '_Tuple2', _0: 'ansi-faint', _1: style.faint}
				]);
			var bgClasses = A3(
				_vito$elm_ansi$Ansi_Log$colorClasses,
				'-bg',
				style.bold,
				_elm_lang$core$Basics$not(style.inverted) ? style.background : style.foreground);
			var fgClasses = A3(
				_vito$elm_ansi$Ansi_Log$colorClasses,
				'-fg',
				style.bold,
				_elm_lang$core$Basics$not(style.inverted) ? style.foreground : style.background);
			var fgbgClasses = A2(
				_elm_lang$core$List$map,
				A2(
					_elm_lang$core$Basics$flip,
					F2(
						function (v0, v1) {
							return {ctor: '_Tuple2', _0: v0, _1: v1};
						}),
					true),
				A2(_elm_lang$core$Basics_ops['++'], fgClasses, bgClasses));
			return _elm_lang$html$Html_Attributes$classList(
				A2(_elm_lang$core$Basics_ops['++'], fgbgClasses, ansiClasses));
		}()
		]);
};
var _vito$elm_ansi$Ansi_Log$viewChunk = function (chunk) {
	return A2(
		_elm_lang$html$Html$span,
		_vito$elm_ansi$Ansi_Log$styleAttributes(chunk.style),
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html$text(chunk.text)
			]));
};
var _vito$elm_ansi$Ansi_Log$viewLine = function (_p1) {
	var _p2 = _p1;
	return A2(
		_elm_lang$html$Html$div,
		_elm_lang$core$Native_List.fromArray(
			[]),
		A3(
			_elm_lang$core$List$foldl,
			F2(
				function (c, l) {
					return A2(
						_elm_lang$core$List_ops['::'],
						_vito$elm_ansi$Ansi_Log$viewChunk(c),
						l);
				}),
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html$text('\n')
				]),
			_p2._0));
};
var _vito$elm_ansi$Ansi_Log$lazyLine = _elm_lang$html$Html_Lazy$lazy(_vito$elm_ansi$Ansi_Log$viewLine);
var _vito$elm_ansi$Ansi_Log$view = function (model) {
	return A2(
		_elm_lang$html$Html$pre,
		_elm_lang$core$Native_List.fromArray(
			[]),
		_elm_lang$core$Array$toList(
			A2(_elm_lang$core$Array$map, _vito$elm_ansi$Ansi_Log$lazyLine, model.lines)));
};
var _vito$elm_ansi$Ansi_Log$blankLine = {
	ctor: '_Tuple2',
	_0: _elm_lang$core$Native_List.fromArray(
		[]),
	_1: 0
};
var _vito$elm_ansi$Ansi_Log$spacing = F2(
	function (style, len) {
		return {
			style: style,
			text: A2(_elm_lang$core$String$repeat, len, ' ')
		};
	});
var _vito$elm_ansi$Ansi_Log$chunkLen = function (_p3) {
	return _elm_lang$core$String$length(
		function (_) {
			return _.text;
		}(_p3));
};
var _vito$elm_ansi$Ansi_Log$addChunk = F2(
	function (chunk, line) {
		var clen = _vito$elm_ansi$Ansi_Log$chunkLen(chunk);
		if (_elm_lang$core$Native_Utils.eq(clen, 0)) {
			return line;
		} else {
			var _p4 = line;
			if (_p4._0.ctor === '[]') {
				return {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Native_List.fromArray(
						[chunk]),
					_1: clen
				};
			} else {
				var _p7 = _p4._1;
				var _p6 = _p4._0._1;
				var _p5 = _p4._0._0;
				return _elm_lang$core$Native_Utils.eq(_p5.style, chunk.style) ? {
					ctor: '_Tuple2',
					_0: A2(
						_elm_lang$core$List_ops['::'],
						_elm_lang$core$Native_Utils.update(
							_p5,
							{
								text: A2(_elm_lang$core$String$append, _p5.text, chunk.text)
							}),
						_p6),
					_1: _p7 + clen
				} : {
					ctor: '_Tuple2',
					_0: A2(
						_elm_lang$core$List_ops['::'],
						chunk,
						A2(_elm_lang$core$List_ops['::'], _p5, _p6)),
					_1: _p7 + clen
				};
			}
		}
	});
var _vito$elm_ansi$Ansi_Log$dropRight = F2(
	function (n, line) {
		dropRight:
		while (true) {
			var _p8 = line;
			if (_p8._0.ctor === '[]') {
				return line;
			} else {
				var _p11 = _p8._1;
				var _p10 = _p8._0._1;
				var _p9 = _p8._0._0;
				var clen = _vito$elm_ansi$Ansi_Log$chunkLen(_p9);
				if (_elm_lang$core$Native_Utils.cmp(clen, n) < 1) {
					var _v4 = n - clen,
						_v5 = {ctor: '_Tuple2', _0: _p10, _1: _p11 - clen};
					n = _v4;
					line = _v5;
					continue dropRight;
				} else {
					return {
						ctor: '_Tuple2',
						_0: A2(
							_elm_lang$core$List_ops['::'],
							_elm_lang$core$Native_Utils.update(
								_p9,
								{
									text: A2(_elm_lang$core$String$dropRight, n, _p9.text)
								}),
							_p10),
						_1: _p11 - n
					};
				}
			}
		}
	});
var _vito$elm_ansi$Ansi_Log$takeRight = F2(
	function (n, line) {
		var _p12 = line;
		if (_p12._0.ctor === '[]') {
			return line;
		} else {
			var _p13 = _p12._0._0;
			var clen = _vito$elm_ansi$Ansi_Log$chunkLen(_p13);
			return (_elm_lang$core$Native_Utils.cmp(clen, n) < 0) ? A2(
				_vito$elm_ansi$Ansi_Log$addChunk,
				_p13,
				A2(
					_vito$elm_ansi$Ansi_Log$takeRight,
					n - clen,
					{ctor: '_Tuple2', _0: _p12._0._1, _1: _p12._1 - clen})) : (_elm_lang$core$Native_Utils.eq(clen, n) ? {
				ctor: '_Tuple2',
				_0: _elm_lang$core$Native_List.fromArray(
					[_p13]),
				_1: clen
			} : {
				ctor: '_Tuple2',
				_0: _elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$core$Native_Utils.update(
						_p13,
						{
							text: A2(_elm_lang$core$String$right, n, _p13.text)
						})
					]),
				_1: n
			});
		}
	});
var _vito$elm_ansi$Ansi_Log$lineLen = _elm_lang$core$Basics$snd;
var _vito$elm_ansi$Ansi_Log$writeChunk = F3(
	function (pos, chunk, line) {
		var len = _vito$elm_ansi$Ansi_Log$lineLen(line);
		var afterLen = len - (_vito$elm_ansi$Ansi_Log$chunkLen(chunk) + pos);
		var textChopped = len - pos;
		if (_elm_lang$core$Native_Utils.eq(len, pos)) {
			return A2(_vito$elm_ansi$Ansi_Log$addChunk, chunk, line);
		} else {
			if (_elm_lang$core$Native_Utils.cmp(pos, len) > 0) {
				return A2(
					_vito$elm_ansi$Ansi_Log$addChunk,
					chunk,
					A2(
						_vito$elm_ansi$Ansi_Log$addChunk,
						A2(_vito$elm_ansi$Ansi_Log$spacing, chunk.style, pos - len),
						line));
			} else {
				var appended = A2(
					_vito$elm_ansi$Ansi_Log$addChunk,
					chunk,
					A2(_vito$elm_ansi$Ansi_Log$dropRight, len - pos, line));
				return (_elm_lang$core$Native_Utils.cmp(afterLen, 0) > 0) ? A3(
					_elm_lang$core$List$foldl,
					_vito$elm_ansi$Ansi_Log$addChunk,
					appended,
					_elm_lang$core$Basics$fst(
						A2(_vito$elm_ansi$Ansi_Log$takeRight, afterLen, line))) : appended;
			}
		}
	});
var _vito$elm_ansi$Ansi_Log$takeLeft = F2(
	function (n, line) {
		return A2(
			_vito$elm_ansi$Ansi_Log$dropRight,
			_vito$elm_ansi$Ansi_Log$lineLen(line) - n,
			line);
	});
var _vito$elm_ansi$Ansi_Log$updateStyle = F2(
	function (action, style) {
		var _p14 = action;
		switch (_p14.ctor) {
			case 'SetForeground':
				return _elm_lang$core$Native_Utils.update(
					style,
					{foreground: _p14._0});
			case 'SetBackground':
				return _elm_lang$core$Native_Utils.update(
					style,
					{background: _p14._0});
			case 'SetInverted':
				return _elm_lang$core$Native_Utils.update(
					style,
					{inverted: _p14._0});
			case 'SetBold':
				return _elm_lang$core$Native_Utils.update(
					style,
					{bold: _p14._0});
			case 'SetFaint':
				return _elm_lang$core$Native_Utils.update(
					style,
					{faint: _p14._0});
			case 'SetItalic':
				return _elm_lang$core$Native_Utils.update(
					style,
					{italic: _p14._0});
			case 'SetUnderline':
				return _elm_lang$core$Native_Utils.update(
					style,
					{underline: _p14._0});
			case 'SetBlink':
				return _elm_lang$core$Native_Utils.update(
					style,
					{blink: _p14._0});
			default:
				return style;
		}
	});
var _vito$elm_ansi$Ansi_Log$appendLine = F3(
	function (after, line, lines) {
		appendLine:
		while (true) {
			if (_elm_lang$core$Native_Utils.eq(after, 0)) {
				return A2(_elm_lang$core$Array$push, line, lines);
			} else {
				var _v8 = after - 1,
					_v9 = line,
					_v10 = A2(_elm_lang$core$Array$push, _vito$elm_ansi$Ansi_Log$blankLine, lines);
				after = _v8;
				line = _v9;
				lines = _v10;
				continue appendLine;
			}
		}
	});
var _vito$elm_ansi$Ansi_Log$updateLine = F3(
	function (row, update, lines) {
		var line = update(
			A2(
				_elm_lang$core$Maybe$withDefault,
				_vito$elm_ansi$Ansi_Log$blankLine,
				A2(_elm_lang$core$Array$get, row, lines)));
		var currentLines = _elm_lang$core$Array$length(lines);
		return (_elm_lang$core$Native_Utils.cmp(row + 1, currentLines) > 0) ? A3(_vito$elm_ansi$Ansi_Log$appendLine, row - currentLines, line, lines) : A3(_elm_lang$core$Array$set, row, line, lines);
	});
var _vito$elm_ansi$Ansi_Log$moveCursor = F3(
	function (r, c, pos) {
		return _elm_lang$core$Native_Utils.update(
			pos,
			{row: pos.row + r, column: pos.column + c});
	});
var _vito$elm_ansi$Ansi_Log$init = function (ldisc) {
	return {
		lineDiscipline: ldisc,
		lines: _elm_lang$core$Array$empty,
		position: {row: 0, column: 0},
		savedPosition: _elm_lang$core$Maybe$Nothing,
		style: {foreground: _elm_lang$core$Maybe$Nothing, background: _elm_lang$core$Maybe$Nothing, bold: false, faint: false, italic: false, underline: false, blink: false, inverted: false},
		remainder: ''
	};
};
var _vito$elm_ansi$Ansi_Log$Model = F6(
	function (a, b, c, d, e, f) {
		return {lineDiscipline: a, lines: b, position: c, savedPosition: d, style: e, remainder: f};
	});
var _vito$elm_ansi$Ansi_Log$Chunk = F2(
	function (a, b) {
		return {text: a, style: b};
	});
var _vito$elm_ansi$Ansi_Log$Style = F8(
	function (a, b, c, d, e, f, g, h) {
		return {foreground: a, background: b, bold: c, faint: d, italic: e, underline: f, blink: g, inverted: h};
	});
var _vito$elm_ansi$Ansi_Log$CursorPosition = F2(
	function (a, b) {
		return {row: a, column: b};
	});
var _vito$elm_ansi$Ansi_Log$handleAction = F2(
	function (action, model) {
		handleAction:
		while (true) {
			var _p15 = action;
			switch (_p15.ctor) {
				case 'Print':
					var chunk = A2(_vito$elm_ansi$Ansi_Log$Chunk, _p15._0, model.style);
					var update = A2(_vito$elm_ansi$Ansi_Log$writeChunk, model.position.column, chunk);
					return _elm_lang$core$Native_Utils.update(
						model,
						{
							lines: A3(_vito$elm_ansi$Ansi_Log$updateLine, model.position.row, update, model.lines),
							position: A3(
								_vito$elm_ansi$Ansi_Log$moveCursor,
								0,
								_vito$elm_ansi$Ansi_Log$chunkLen(chunk),
								model.position)
						});
				case 'CarriageReturn':
					return _elm_lang$core$Native_Utils.update(
						model,
						{
							position: A2(_vito$elm_ansi$Ansi_Log$CursorPosition, model.position.row, 0)
						});
				case 'Linebreak':
					var _v13 = _vito$elm_ansi$Ansi$Print(''),
						_v14 = function () {
						var _p16 = model.lineDiscipline;
						if (_p16.ctor === 'Raw') {
							return _elm_lang$core$Native_Utils.update(
								model,
								{
									position: A3(_vito$elm_ansi$Ansi_Log$moveCursor, 1, 0, model.position)
								});
						} else {
							return _elm_lang$core$Native_Utils.update(
								model,
								{
									position: A2(_vito$elm_ansi$Ansi_Log$CursorPosition, model.position.row + 1, 0)
								});
						}
					}();
					action = _v13;
					model = _v14;
					continue handleAction;
				case 'Remainder':
					return _elm_lang$core$Native_Utils.update(
						model,
						{remainder: _p15._0});
				case 'CursorUp':
					return _elm_lang$core$Native_Utils.update(
						model,
						{
							position: A3(_vito$elm_ansi$Ansi_Log$moveCursor, 0 - _p15._0, 0, model.position)
						});
				case 'CursorDown':
					return _elm_lang$core$Native_Utils.update(
						model,
						{
							position: A3(_vito$elm_ansi$Ansi_Log$moveCursor, _p15._0, 0, model.position)
						});
				case 'CursorForward':
					return _elm_lang$core$Native_Utils.update(
						model,
						{
							position: A3(_vito$elm_ansi$Ansi_Log$moveCursor, 0, _p15._0, model.position)
						});
				case 'CursorBack':
					return _elm_lang$core$Native_Utils.update(
						model,
						{
							position: A3(_vito$elm_ansi$Ansi_Log$moveCursor, 0, 0 - _p15._0, model.position)
						});
				case 'CursorPosition':
					return _elm_lang$core$Native_Utils.update(
						model,
						{
							position: A2(_vito$elm_ansi$Ansi_Log$CursorPosition, _p15._0 - 1, _p15._1 - 1)
						});
				case 'CursorColumn':
					return _elm_lang$core$Native_Utils.update(
						model,
						{
							position: A2(_vito$elm_ansi$Ansi_Log$CursorPosition, model.position.row, _p15._0)
						});
				case 'SaveCursorPosition':
					return _elm_lang$core$Native_Utils.update(
						model,
						{
							savedPosition: _elm_lang$core$Maybe$Just(model.position)
						});
				case 'RestoreCursorPosition':
					return _elm_lang$core$Native_Utils.update(
						model,
						{
							position: A2(_elm_lang$core$Maybe$withDefault, model.position, model.savedPosition)
						});
				case 'EraseLine':
					var _p17 = _p15._0;
					switch (_p17.ctor) {
						case 'EraseToBeginning':
							var chunk = A2(
								_vito$elm_ansi$Ansi_Log$Chunk,
								A2(_elm_lang$core$String$repeat, model.position.column, ' '),
								model.style);
							var update = A2(_vito$elm_ansi$Ansi_Log$writeChunk, 0, chunk);
							return _elm_lang$core$Native_Utils.update(
								model,
								{
									lines: A3(_vito$elm_ansi$Ansi_Log$updateLine, model.position.row, update, model.lines)
								});
						case 'EraseToEnd':
							var update = _vito$elm_ansi$Ansi_Log$takeLeft(model.position.column);
							return _elm_lang$core$Native_Utils.update(
								model,
								{
									lines: A3(_vito$elm_ansi$Ansi_Log$updateLine, model.position.row, update, model.lines)
								});
						default:
							return _elm_lang$core$Native_Utils.update(
								model,
								{
									lines: A3(
										_vito$elm_ansi$Ansi_Log$updateLine,
										model.position.row,
										_elm_lang$core$Basics$always(_vito$elm_ansi$Ansi_Log$blankLine),
										model.lines)
								});
					}
				default:
					return _elm_lang$core$Native_Utils.update(
						model,
						{
							style: A2(_vito$elm_ansi$Ansi_Log$updateStyle, action, model.style)
						});
			}
		}
	});
var _vito$elm_ansi$Ansi_Log$update = F2(
	function (str, model) {
		return A3(
			_vito$elm_ansi$Ansi$parseInto,
			_elm_lang$core$Native_Utils.update(
				model,
				{remainder: ''}),
			_vito$elm_ansi$Ansi_Log$handleAction,
			A2(_elm_lang$core$Basics_ops['++'], model.remainder, str));
	});
var _vito$elm_ansi$Ansi_Log$Cooked = {ctor: 'Cooked'};
var _vito$elm_ansi$Ansi_Log$Raw = {ctor: 'Raw'};

var _concourse$atc$Concourse_Pagination$parseNum = function (_p0) {
	return _elm_lang$core$Result$toMaybe(
		_elm_lang$core$String$toInt(_p0));
};
var _concourse$atc$Concourse_Pagination$toQuery = function (page) {
	var _p1 = page;
	if (_p1.ctor === 'Nothing') {
		return _elm_lang$core$Dict$empty;
	} else {
		var limitParam = {
			ctor: '_Tuple2',
			_0: 'limit',
			_1: _elm_lang$core$Basics$toString(_p1._0.limit)
		};
		var directionParam = function () {
			var _p2 = _p1._0.direction;
			if (_p2.ctor === 'Since') {
				return {
					ctor: '_Tuple2',
					_0: 'since',
					_1: _elm_lang$core$Basics$toString(_p2._0)
				};
			} else {
				return {
					ctor: '_Tuple2',
					_0: 'until',
					_1: _elm_lang$core$Basics$toString(_p2._0)
				};
			}
		}();
		return _elm_lang$core$Dict$fromList(
			_elm_lang$core$Native_List.fromArray(
				[directionParam, limitParam]));
	}
};
var _concourse$atc$Concourse_Pagination$parseQuery = function (query) {
	var parseParam = function (p) {
		var _p3 = A2(_elm_lang$core$String$split, '=', p);
		if (_p3.ctor === '::') {
			return {
				ctor: '_Tuple2',
				_0: _p3._0,
				_1: A2(_elm_lang$core$String$join, '=', _p3._1)
			};
		} else {
			return {ctor: '_Tuple2', _0: '', _1: ''};
		}
	};
	return _elm_lang$core$Dict$fromList(
		A2(
			_elm_lang$core$List$map,
			parseParam,
			A2(_elm_lang$core$String$split, '&', query)));
};
var _concourse$atc$Concourse_Pagination$setQuery = F2(
	function (baseURL, query) {
		var params = A2(
			_elm_lang$core$String$join,
			'&',
			A2(
				_elm_lang$core$List$map,
				function (_p4) {
					var _p5 = _p4;
					return A2(
						_elm_lang$core$Basics_ops['++'],
						_p5._0,
						A2(_elm_lang$core$Basics_ops['++'], '=', _p5._1));
				},
				_elm_lang$core$Dict$toList(query)));
		return _elm_lang$core$Native_Utils.eq(params, '') ? baseURL : A2(
			_elm_lang$core$Basics_ops['++'],
			baseURL,
			A2(_elm_lang$core$Basics_ops['++'], '?', params));
	});
var _concourse$atc$Concourse_Pagination$extractQuery = function (url) {
	var _p6 = A2(_elm_lang$core$String$split, '?', url);
	if ((_p6.ctor === '::') && (_p6._1.ctor === '::')) {
		return {
			ctor: '_Tuple2',
			_0: _p6._0,
			_1: _concourse$atc$Concourse_Pagination$parseQuery(_p6._1._0)
		};
	} else {
		return {ctor: '_Tuple2', _0: url, _1: _elm_lang$core$Dict$empty};
	}
};
var _concourse$atc$Concourse_Pagination$addParams = F2(
	function (url, page) {
		var _p7 = _concourse$atc$Concourse_Pagination$extractQuery(url);
		var baseURL = _p7._0;
		var query = _p7._1;
		return A2(
			_concourse$atc$Concourse_Pagination$setQuery,
			baseURL,
			A2(
				_elm_lang$core$Dict$union,
				query,
				_concourse$atc$Concourse_Pagination$toQuery(page)));
	});
var _concourse$atc$Concourse_Pagination$fstToLower = function (_p8) {
	var _p9 = _p8;
	return {
		ctor: '_Tuple2',
		_0: _elm_lang$core$String$toLower(_p9._0),
		_1: _p9._1
	};
};
var _concourse$atc$Concourse_Pagination$keysToLower = function (_p10) {
	return _elm_lang$core$Dict$fromList(
		A2(
			_elm_lang$core$List$map,
			_concourse$atc$Concourse_Pagination$fstToLower,
			_elm_lang$core$Dict$toList(_p10)));
};
var _concourse$atc$Concourse_Pagination$promoteHttpError = function (rawError) {
	var _p11 = rawError;
	if (_p11.ctor === 'RawTimeout') {
		return _evancz$elm_http$Http$Timeout;
	} else {
		return _evancz$elm_http$Http$NetworkError;
	}
};
var _concourse$atc$Concourse_Pagination$handleResponse = function (response) {
	if ((_elm_lang$core$Native_Utils.cmp(200, response.status) < 1) && (_elm_lang$core$Native_Utils.cmp(response.status, 300) < 0)) {
		var _p12 = response.value;
		if (_p12.ctor === 'Text') {
			return _elm_lang$core$Result$Ok(_p12._0);
		} else {
			return _elm_lang$core$Result$Err(
				_evancz$elm_http$Http$UnexpectedPayload('Response body is a blob, expecting a string.'));
		}
	} else {
		return _elm_lang$core$Result$Err(
			A2(_evancz$elm_http$Http$BadResponse, response.status, response.statusText));
	}
};
var _concourse$atc$Concourse_Pagination$nextRel = 'next';
var _concourse$atc$Concourse_Pagination$previousRel = 'previous';
var _concourse$atc$Concourse_Pagination$linkHeaderRegex = _elm_lang$core$Regex$regex(
	A2(
		_elm_lang$core$Basics_ops['++'],
		'<([^>]+)>; rel=\"(',
		A2(
			_elm_lang$core$Basics_ops['++'],
			_concourse$atc$Concourse_Pagination$previousRel,
			A2(
				_elm_lang$core$Basics_ops['++'],
				'|',
				A2(_elm_lang$core$Basics_ops['++'], _concourse$atc$Concourse_Pagination$nextRel, ')\"')))));
var _concourse$atc$Concourse_Pagination$parseLinkTuple = function (header) {
	var _p13 = A3(
		_elm_lang$core$Regex$find,
		_elm_lang$core$Regex$AtMost(1),
		_concourse$atc$Concourse_Pagination$linkHeaderRegex,
		header);
	if (_p13.ctor === '[]') {
		return _elm_lang$core$Maybe$Nothing;
	} else {
		var _p14 = _p13._0.submatches;
		if ((((_p14.ctor === '::') && (_p14._0.ctor === 'Just')) && (_p14._1.ctor === '::')) && (_p14._1._0.ctor === 'Just')) {
			return _elm_lang$core$Maybe$Just(
				{ctor: '_Tuple2', _0: _p14._1._0._0, _1: _p14._0._0});
		} else {
			return _elm_lang$core$Maybe$Nothing;
		}
	}
};
var _concourse$atc$Concourse_Pagination$Paginated = F2(
	function (a, b) {
		return {content: a, pagination: b};
	});
var _concourse$atc$Concourse_Pagination$Pagination = F2(
	function (a, b) {
		return {previousPage: a, nextPage: b};
	});
var _concourse$atc$Concourse_Pagination$Page = F2(
	function (a, b) {
		return {direction: a, limit: b};
	});
var _concourse$atc$Concourse_Pagination$Until = function (a) {
	return {ctor: 'Until', _0: a};
};
var _concourse$atc$Concourse_Pagination$Since = function (a) {
	return {ctor: 'Since', _0: a};
};
var _concourse$atc$Concourse_Pagination$fromQuery = function (query) {
	var since = A2(
		_elm_lang$core$Maybe$map,
		_concourse$atc$Concourse_Pagination$Since,
		A2(
			_elm_lang$core$Maybe$andThen,
			A2(_elm_lang$core$Dict$get, 'since', query),
			_concourse$atc$Concourse_Pagination$parseNum));
	var until = A2(
		_elm_lang$core$Maybe$map,
		_concourse$atc$Concourse_Pagination$Until,
		A2(
			_elm_lang$core$Maybe$andThen,
			A2(_elm_lang$core$Dict$get, 'until', query),
			_concourse$atc$Concourse_Pagination$parseNum));
	var limit = A2(
		_elm_lang$core$Maybe$withDefault,
		0,
		A2(
			_elm_lang$core$Maybe$andThen,
			A2(_elm_lang$core$Dict$get, 'limit', query),
			_concourse$atc$Concourse_Pagination$parseNum));
	return A2(
		_elm_lang$core$Maybe$map,
		function (direction) {
			return {direction: direction, limit: limit};
		},
		_elm_lang$core$Maybe$oneOf(
			_elm_lang$core$Native_List.fromArray(
				[until, since])));
};
var _concourse$atc$Concourse_Pagination$parseParams = function (_p15) {
	return _concourse$atc$Concourse_Pagination$fromQuery(
		_elm_lang$core$Basics$snd(
			_concourse$atc$Concourse_Pagination$extractQuery(_p15)));
};
var _concourse$atc$Concourse_Pagination$parseLinks = function (response) {
	var _p16 = A2(
		_elm_lang$core$Dict$get,
		'link',
		_concourse$atc$Concourse_Pagination$keysToLower(response.headers));
	if (_p16.ctor === 'Nothing') {
		return A2(_concourse$atc$Concourse_Pagination$Pagination, _elm_lang$core$Maybe$Nothing, _elm_lang$core$Maybe$Nothing);
	} else {
		var headers = A2(_elm_lang$core$String$split, ', ', _p16._0);
		var parsed = _elm_lang$core$Dict$fromList(
			A2(_elm_lang$core$List$filterMap, _concourse$atc$Concourse_Pagination$parseLinkTuple, headers));
		return A2(
			_concourse$atc$Concourse_Pagination$Pagination,
			A2(
				_elm_lang$core$Maybe$andThen,
				A2(_elm_lang$core$Dict$get, _concourse$atc$Concourse_Pagination$previousRel, parsed),
				_concourse$atc$Concourse_Pagination$parseParams),
			A2(
				_elm_lang$core$Maybe$andThen,
				A2(_elm_lang$core$Dict$get, _concourse$atc$Concourse_Pagination$nextRel, parsed),
				_concourse$atc$Concourse_Pagination$parseParams));
	}
};
var _concourse$atc$Concourse_Pagination$parsePagination = F2(
	function (decode, response) {
		var decoded = A2(
			_elm_lang$core$Result$andThen,
			_concourse$atc$Concourse_Pagination$handleResponse(response),
			function (body) {
				return A2(
					_elm_lang$core$Result$formatError,
					_evancz$elm_http$Http$UnexpectedPayload,
					A2(
						_elm_lang$core$Json_Decode$decodeString,
						_elm_lang$core$Json_Decode$list(decode),
						body));
			});
		var pagination = _concourse$atc$Concourse_Pagination$parseLinks(response);
		var _p17 = decoded;
		if (_p17.ctor === 'Err') {
			return _elm_lang$core$Task$fail(_p17._0);
		} else {
			return _elm_lang$core$Task$succeed(
				{content: _p17._0, pagination: pagination});
		}
	});
var _concourse$atc$Concourse_Pagination$fetch = F3(
	function (decode, url, page) {
		var get = A2(
			_evancz$elm_http$Http$send,
			_evancz$elm_http$Http$defaultSettings,
			{
				verb: 'GET',
				headers: _elm_lang$core$Native_List.fromArray(
					[]),
				url: A2(_concourse$atc$Concourse_Pagination$addParams, url, page),
				body: _evancz$elm_http$Http$empty
			});
		return A2(
			_elm_lang$core$Task$andThen,
			A2(_elm_lang$core$Task$mapError, _concourse$atc$Concourse_Pagination$promoteHttpError, get),
			_concourse$atc$Concourse_Pagination$parsePagination(decode));
	});

var _concourse$atc$Concourse_BuildStatus$isRunning = function (status) {
	var _p0 = status;
	switch (_p0.ctor) {
		case 'Pending':
			return true;
		case 'Started':
			return true;
		default:
			return false;
	}
};
var _concourse$atc$Concourse_BuildStatus$show = function (status) {
	var _p1 = status;
	switch (_p1.ctor) {
		case 'Pending':
			return 'pending';
		case 'Started':
			return 'started';
		case 'Succeeded':
			return 'succeeded';
		case 'Failed':
			return 'failed';
		case 'Errored':
			return 'errored';
		default:
			return 'aborted';
	}
};
var _concourse$atc$Concourse_BuildStatus$Aborted = {ctor: 'Aborted'};
var _concourse$atc$Concourse_BuildStatus$Errored = {ctor: 'Errored'};
var _concourse$atc$Concourse_BuildStatus$Failed = {ctor: 'Failed'};
var _concourse$atc$Concourse_BuildStatus$Succeeded = {ctor: 'Succeeded'};
var _concourse$atc$Concourse_BuildStatus$Started = {ctor: 'Started'};
var _concourse$atc$Concourse_BuildStatus$Pending = {ctor: 'Pending'};
var _concourse$atc$Concourse_BuildStatus$decode = A2(
	_elm_lang$core$Json_Decode$customDecoder,
	_elm_lang$core$Json_Decode$string,
	function (status) {
		var _p2 = status;
		switch (_p2) {
			case 'pending':
				return _elm_lang$core$Result$Ok(_concourse$atc$Concourse_BuildStatus$Pending);
			case 'started':
				return _elm_lang$core$Result$Ok(_concourse$atc$Concourse_BuildStatus$Started);
			case 'succeeded':
				return _elm_lang$core$Result$Ok(_concourse$atc$Concourse_BuildStatus$Succeeded);
			case 'failed':
				return _elm_lang$core$Result$Ok(_concourse$atc$Concourse_BuildStatus$Failed);
			case 'errored':
				return _elm_lang$core$Result$Ok(_concourse$atc$Concourse_BuildStatus$Errored);
			case 'aborted':
				return _elm_lang$core$Result$Ok(_concourse$atc$Concourse_BuildStatus$Aborted);
			default:
				return _elm_lang$core$Result$Err(
					A2(_elm_lang$core$Basics_ops['++'], 'unknown build status: ', _p2));
		}
	});

var _concourse$atc$Concourse_Build$dateFromSeconds = function (_p0) {
	return _elm_lang$core$Date$fromTime(
		A2(
			F2(
				function (x, y) {
					return x * y;
				}),
			1000,
			_p0));
};
var _concourse$atc$Concourse_Build$promoteHttpError = function (rawError) {
	var _p1 = rawError;
	if (_p1.ctor === 'RawTimeout') {
		return _evancz$elm_http$Http$Timeout;
	} else {
		return _evancz$elm_http$Http$NetworkError;
	}
};
var _concourse$atc$Concourse_Build$handleResponse = function (response) {
	return ((_elm_lang$core$Native_Utils.cmp(200, response.status) < 1) && (_elm_lang$core$Native_Utils.cmp(response.status, 300) < 0)) ? _elm_lang$core$Task$succeed(
		{ctor: '_Tuple0'}) : _elm_lang$core$Task$fail(
		A2(_evancz$elm_http$Http$BadResponse, response.status, response.statusText));
};
var _concourse$atc$Concourse_Build$url = function (build) {
	var _p2 = build.job;
	if (_p2.ctor === 'Nothing') {
		return A2(
			_elm_lang$core$Basics_ops['++'],
			'/builds/',
			_elm_lang$core$Basics$toString(build.id));
	} else {
		return A2(
			_elm_lang$core$Basics_ops['++'],
			'/teams/',
			A2(
				_elm_lang$core$Basics_ops['++'],
				_p2._0.teamName,
				A2(
					_elm_lang$core$Basics_ops['++'],
					'/pipelines/',
					A2(
						_elm_lang$core$Basics_ops['++'],
						_p2._0.pipelineName,
						A2(
							_elm_lang$core$Basics_ops['++'],
							'/jobs/',
							A2(
								_elm_lang$core$Basics_ops['++'],
								_p2._0.name,
								A2(_elm_lang$core$Basics_ops['++'], '/builds/', build.name)))))));
	}
};
var _concourse$atc$Concourse_Build$abort = function (buildId) {
	var post = A2(
		_evancz$elm_http$Http$send,
		_evancz$elm_http$Http$defaultSettings,
		{
			verb: 'POST',
			headers: _elm_lang$core$Native_List.fromArray(
				[]),
			url: A2(
				_elm_lang$core$Basics_ops['++'],
				'/api/v1/builds/',
				A2(
					_elm_lang$core$Basics_ops['++'],
					_elm_lang$core$Basics$toString(buildId),
					'/abort')),
			body: _evancz$elm_http$Http$empty
		});
	return A2(
		_elm_lang$core$Task$andThen,
		A2(_elm_lang$core$Task$mapError, _concourse$atc$Concourse_Build$promoteHttpError, post),
		_concourse$atc$Concourse_Build$handleResponse);
};
var _concourse$atc$Concourse_Build$Build = F7(
	function (a, b, c, d, e, f, g) {
		return {id: a, url: b, name: c, job: d, status: e, duration: f, reapTime: g};
	});
var _concourse$atc$Concourse_Build$BuildJob = F3(
	function (a, b, c) {
		return {name: a, teamName: b, pipelineName: c};
	});
var _concourse$atc$Concourse_Build$BuildDuration = F2(
	function (a, b) {
		return {startedAt: a, finishedAt: b};
	});
var _concourse$atc$Concourse_Build$decode = A8(
	_elm_lang$core$Json_Decode$object7,
	_concourse$atc$Concourse_Build$Build,
	A2(_elm_lang$core$Json_Decode_ops[':='], 'id', _elm_lang$core$Json_Decode$int),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'url', _elm_lang$core$Json_Decode$string),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'name', _elm_lang$core$Json_Decode$string),
	_elm_lang$core$Json_Decode$maybe(
		A4(
			_elm_lang$core$Json_Decode$object3,
			_concourse$atc$Concourse_Build$BuildJob,
			A2(_elm_lang$core$Json_Decode_ops[':='], 'job_name', _elm_lang$core$Json_Decode$string),
			A2(_elm_lang$core$Json_Decode_ops[':='], 'team_name', _elm_lang$core$Json_Decode$string),
			A2(_elm_lang$core$Json_Decode_ops[':='], 'pipeline_name', _elm_lang$core$Json_Decode$string))),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'status', _concourse$atc$Concourse_BuildStatus$decode),
	A3(
		_elm_lang$core$Json_Decode$object2,
		_concourse$atc$Concourse_Build$BuildDuration,
		_elm_lang$core$Json_Decode$maybe(
			A2(
				_elm_lang$core$Json_Decode_ops[':='],
				'start_time',
				A2(_elm_lang$core$Json_Decode$map, _concourse$atc$Concourse_Build$dateFromSeconds, _elm_lang$core$Json_Decode$float))),
		_elm_lang$core$Json_Decode$maybe(
			A2(
				_elm_lang$core$Json_Decode_ops[':='],
				'end_time',
				A2(_elm_lang$core$Json_Decode$map, _concourse$atc$Concourse_Build$dateFromSeconds, _elm_lang$core$Json_Decode$float)))),
	_elm_lang$core$Json_Decode$maybe(
		A2(
			_elm_lang$core$Json_Decode_ops[':='],
			'reap_time',
			A2(_elm_lang$core$Json_Decode$map, _concourse$atc$Concourse_Build$dateFromSeconds, _elm_lang$core$Json_Decode$float))));
var _concourse$atc$Concourse_Build$fetch = function (buildId) {
	return A2(
		_evancz$elm_http$Http$get,
		_concourse$atc$Concourse_Build$decode,
		A2(
			_elm_lang$core$Basics_ops['++'],
			'/api/v1/builds/',
			_elm_lang$core$Basics$toString(buildId)));
};
var _concourse$atc$Concourse_Build$fetchJobBuilds = F2(
	function (job, page) {
		var url = A2(
			_elm_lang$core$Basics_ops['++'],
			'/api/v1/teams/',
			A2(
				_elm_lang$core$Basics_ops['++'],
				job.teamName,
				A2(
					_elm_lang$core$Basics_ops['++'],
					'/pipelines/',
					A2(
						_elm_lang$core$Basics_ops['++'],
						job.pipelineName,
						A2(
							_elm_lang$core$Basics_ops['++'],
							'/jobs/',
							A2(_elm_lang$core$Basics_ops['++'], job.name, '/builds'))))));
		return A3(_concourse$atc$Concourse_Pagination$fetch, _concourse$atc$Concourse_Build$decode, url, page);
	});

var _concourse$atc$Concourse_Metadata$MetadataField = F2(
	function (a, b) {
		return {name: a, value: b};
	});
var _concourse$atc$Concourse_Metadata$decodeMetadataField = A3(
	_elm_lang$core$Json_Decode$object2,
	_concourse$atc$Concourse_Metadata$MetadataField,
	A2(_elm_lang$core$Json_Decode_ops[':='], 'name', _elm_lang$core$Json_Decode$string),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'value', _elm_lang$core$Json_Decode$string));
var _concourse$atc$Concourse_Metadata$decode = _elm_lang$core$Json_Decode$list(_concourse$atc$Concourse_Metadata$decodeMetadataField);

var _concourse$atc$Concourse_Version$decode = _elm_lang$core$Json_Decode$dict(_elm_lang$core$Json_Decode$string);

var _concourse$atc$Native_EventSource = function() {
  function open(url, settings) {
    return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback) {
      var source = new EventSource(url);

      function dispatchEvent(event) {
        var ev = {
          data: event.data
        }

        if (event.type !== undefined) {
          ev.name = _elm_lang$core$Maybe$Just(event.type);
        } else {
          ev.name = _elm_lang$core$Maybe$Nothing;
        }

        if (event.lastEventId !== undefined) {
          ev.lastEventId = _elm_lang$core$Maybe$Just(event.lastEventId);
        } else {
          ev.lastEventId = _elm_lang$core$Maybe$Nothing;
        }

        _elm_lang$core$Native_Scheduler.rawSpawn(settings.onEvent(ev));
      };

      source.onmessage = function(event) {
        dispatchEvent(event);
      };

      _elm_lang$core$Native_List.toArray(settings.events).forEach(function(eventType) {
        source.addEventListener(eventType, function(event) {
          dispatchEvent(event);
        });
      });

      source.onopen = function(event) {
        _elm_lang$core$Native_Scheduler.rawSpawn(settings.onOpen(source));
      };

      source.onerror = function(event) {
        _elm_lang$core$Native_Scheduler.rawSpawn(settings.onError(_elm_lang$core$Native_Utils.Tuple0));
      };
    });
  }

  function close(source) {
    return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback) {
      source.close();
      callback(_elm_lang$core$Native_Scheduler.succeed(_elm_lang$core$Native_Utils.Tuple0));
    });
  }

  return {
    open: F2(open),
    close: close
  };
}();

var _concourse$atc$EventSource_LowLevel$close = _concourse$atc$Native_EventSource.close;
var _concourse$atc$EventSource_LowLevel$open = _concourse$atc$Native_EventSource.open;
var _concourse$atc$EventSource_LowLevel$Event = F3(
	function (a, b, c) {
		return {lastEventId: a, name: b, data: c};
	});
var _concourse$atc$EventSource_LowLevel$Settings = F4(
	function (a, b, c, d) {
		return {events: a, onEvent: b, onOpen: c, onError: d};
	});
var _concourse$atc$EventSource_LowLevel$EventSource = {ctor: 'EventSource'};

var _concourse$atc$EventSource$broadcast = F3(
	function (router, msg, subs) {
		return A2(
			_elm_lang$core$Task$map,
			_elm_lang$core$Basics$always(
				{ctor: '_Tuple0'}),
			_elm_lang$core$Task$sequence(
				A2(
					_elm_lang$core$List$map,
					function (_p0) {
						var _p1 = _p0;
						return A2(
							_elm_lang$core$Platform$sendToApp,
							router,
							_p1._1(msg));
					},
					subs)));
	});
var _concourse$atc$EventSource$closeSource = F3(
	function (key, source, rest) {
		return A2(
			_elm_lang$core$Task$andThen,
			rest,
			function (state) {
				var _p2 = source.source;
				if (_p2.ctor === 'Nothing') {
					return _elm_lang$core$Task$succeed(
						A2(_elm_lang$core$Dict$remove, key, state));
				} else {
					return A2(
						_elm_lang$core$Task$andThen,
						_concourse$atc$EventSource_LowLevel$close(_p2._0),
						function (_p3) {
							return _elm_lang$core$Task$succeed(
								A2(_elm_lang$core$Dict$remove, key, state));
						});
				}
			});
	});
var _concourse$atc$EventSource$updateSourceSubs = F4(
	function (key, subs, source, rest) {
		return A2(
			_elm_lang$core$Task$map,
			A2(
				_elm_lang$core$Dict$insert,
				key,
				_elm_lang$core$Native_Utils.update(
					source,
					{subs: subs})),
			rest);
	});
var _concourse$atc$EventSource$init = _elm_lang$core$Task$succeed(_elm_lang$core$Dict$empty);
var _concourse$atc$EventSource$subscription = _elm_lang$core$Native_Platform.leaf('EventSource');
var _concourse$atc$EventSource$Source = F3(
	function (a, b, c) {
		return {subs: a, watcher: b, source: c};
	});
var _concourse$atc$EventSource$ESErrored = function (a) {
	return {ctor: 'ESErrored', _0: a};
};
var _concourse$atc$EventSource$ESOpened = F2(
	function (a, b) {
		return {ctor: 'ESOpened', _0: a, _1: b};
	});
var _concourse$atc$EventSource$ESEvent = F2(
	function (a, b) {
		return {ctor: 'ESEvent', _0: a, _1: b};
	});
var _concourse$atc$EventSource$open = F2(
	function (router, _p4) {
		var _p5 = _p4;
		var _p10 = _p5._0;
		var _p9 = _p5._1;
		return A2(
			_concourse$atc$EventSource_LowLevel$open,
			_p10,
			{
				events: _p9,
				onEvent: function (_p6) {
					return A2(
						_elm_lang$core$Platform$sendToSelf,
						router,
						A2(
							_concourse$atc$EventSource$ESEvent,
							{ctor: '_Tuple2', _0: _p10, _1: _p9},
							_p6));
				},
				onOpen: function (_p7) {
					return A2(
						_elm_lang$core$Platform$sendToSelf,
						router,
						A2(
							_concourse$atc$EventSource$ESOpened,
							{ctor: '_Tuple2', _0: _p10, _1: _p9},
							_p7));
				},
				onError: function (_p8) {
					return A2(
						_elm_lang$core$Platform$sendToSelf,
						router,
						A2(
							_elm_lang$core$Basics$always,
							_concourse$atc$EventSource$ESErrored(
								{ctor: '_Tuple2', _0: _p10, _1: _p9}),
							_p8));
				}
			});
	});
var _concourse$atc$EventSource$createSource = F4(
	function (router, key, subs, rest) {
		return A2(
			_elm_lang$core$Task$andThen,
			rest,
			function (state) {
				return A2(
					_elm_lang$core$Task$andThen,
					_elm_lang$core$Process$spawn(
						A2(_concourse$atc$EventSource$open, router, key)),
					function (processId) {
						return _elm_lang$core$Task$succeed(
							A3(
								_elm_lang$core$Dict$insert,
								key,
								{subs: subs, watcher: processId, source: _elm_lang$core$Maybe$Nothing},
								state));
					});
			});
	});
var _concourse$atc$EventSource$Errored = {ctor: 'Errored'};
var _concourse$atc$EventSource$Opened = {ctor: 'Opened'};
var _concourse$atc$EventSource$Event = function (a) {
	return {ctor: 'Event', _0: a};
};
var _concourse$atc$EventSource$onSelfMsg = F3(
	function (router, msg, state) {
		var _p11 = msg;
		switch (_p11.ctor) {
			case 'ESEvent':
				var _p12 = A2(_elm_lang$core$Dict$get, _p11._0, state);
				if (_p12.ctor === 'Nothing') {
					return _elm_lang$core$Task$succeed(state);
				} else {
					return A2(
						_elm_lang$core$Task$andThen,
						A3(
							_concourse$atc$EventSource$broadcast,
							router,
							_concourse$atc$EventSource$Event(_p11._1),
							_p12._0.subs),
						function (_p13) {
							return _elm_lang$core$Task$succeed(state);
						});
				}
			case 'ESOpened':
				var _p17 = _p11._0;
				var _p14 = A2(_elm_lang$core$Dict$get, _p17, state);
				if (_p14.ctor === 'Nothing') {
					return _elm_lang$core$Task$succeed(state);
				} else {
					var _p16 = _p14._0;
					return A2(
						_elm_lang$core$Task$andThen,
						A3(_concourse$atc$EventSource$broadcast, router, _concourse$atc$EventSource$Opened, _p16.subs),
						function (_p15) {
							return _elm_lang$core$Task$succeed(
								A3(
									_elm_lang$core$Dict$insert,
									_p17,
									_elm_lang$core$Native_Utils.update(
										_p16,
										{
											source: _elm_lang$core$Maybe$Just(_p11._1)
										}),
									state));
						});
				}
			default:
				var _p18 = A2(_elm_lang$core$Dict$get, _p11._0, state);
				if (_p18.ctor === 'Nothing') {
					return _elm_lang$core$Task$succeed(state);
				} else {
					return A2(
						_elm_lang$core$Task$andThen,
						A3(_concourse$atc$EventSource$broadcast, router, _concourse$atc$EventSource$Errored, _p18._0.subs),
						function (_p19) {
							return _elm_lang$core$Task$succeed(state);
						});
				}
		}
	});
var _concourse$atc$EventSource$MySub = F2(
	function (a, b) {
		return {ctor: 'MySub', _0: a, _1: b};
	});
var _concourse$atc$EventSource$listen = F2(
	function (key, tagger) {
		return _concourse$atc$EventSource$subscription(
			A2(_concourse$atc$EventSource$MySub, key, tagger));
	});
var _concourse$atc$EventSource$subMap = F2(
	function (func, _p20) {
		var _p21 = _p20;
		return A2(
			_concourse$atc$EventSource$MySub,
			_p21._0,
			function (_p22) {
				return func(
					_p21._1(_p22));
			});
	});
var _concourse$atc$EventSource$onEffects = F3(
	function (router, subs, state) {
		var addSub = F3(
			function (key, tagger, msubs) {
				return _elm_lang$core$Maybe$Just(
					A2(
						_elm_lang$core$List_ops['::'],
						A2(_concourse$atc$EventSource$MySub, key, tagger),
						A2(
							_elm_lang$core$Maybe$withDefault,
							_elm_lang$core$Native_List.fromArray(
								[]),
							msubs)));
			});
		var insertSub = F2(
			function (_p23, state) {
				var _p24 = _p23;
				var _p25 = _p24._0;
				return A3(
					_elm_lang$core$Dict$update,
					_p25,
					A2(addSub, _p25, _p24._1),
					state);
			});
		var desiredSubs = A3(_elm_lang$core$List$foldl, insertSub, _elm_lang$core$Dict$empty, subs);
		return A6(
			_elm_lang$core$Dict$merge,
			_concourse$atc$EventSource$createSource(router),
			_concourse$atc$EventSource$updateSourceSubs,
			_concourse$atc$EventSource$closeSource,
			desiredSubs,
			state,
			_elm_lang$core$Task$succeed(_elm_lang$core$Dict$empty));
	});
_elm_lang$core$Native_Platform.effectManagers['EventSource'] = {pkg: 'concourse/atc', init: _concourse$atc$EventSource$init, onEffects: _concourse$atc$EventSource$onEffects, onSelfMsg: _concourse$atc$EventSource$onSelfMsg, tag: 'sub', subMap: _concourse$atc$EventSource$subMap};

var _concourse$atc$Concourse_BuildEvents$dateFromSeconds = function (_p0) {
	return _elm_lang$core$Date$fromTime(
		A2(
			F2(
				function (x, y) {
					return x * y;
				}),
			1000,
			_p0));
};
var _concourse$atc$Concourse_BuildEvents$BuildEventEnvelope = F3(
	function (a, b, c) {
		return {event: a, version: b, value: c};
	});
var _concourse$atc$Concourse_BuildEvents$decodeEnvelope = A4(
	_elm_lang$core$Json_Decode$object3,
	_concourse$atc$Concourse_BuildEvents$BuildEventEnvelope,
	A2(_elm_lang$core$Json_Decode_ops[':='], 'event', _elm_lang$core$Json_Decode$string),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'version', _elm_lang$core$Json_Decode$string),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'data', _elm_lang$core$Json_Decode$value));
var _concourse$atc$Concourse_BuildEvents$Origin = F2(
	function (a, b) {
		return {source: a, id: b};
	});
var _concourse$atc$Concourse_BuildEvents$decodeOrigin = A3(
	_elm_lang$core$Json_Decode$object2,
	_concourse$atc$Concourse_BuildEvents$Origin,
	function (_p1) {
		return A2(
			_elm_lang$core$Json_Decode$map,
			_elm_lang$core$Maybe$withDefault(''),
			_elm_lang$core$Json_Decode$maybe(_p1));
	}(
		A2(_elm_lang$core$Json_Decode_ops[':='], 'source', _elm_lang$core$Json_Decode$string)),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'id', _elm_lang$core$Json_Decode$string));
var _concourse$atc$Concourse_BuildEvents$decodeFinishResource = function (cons) {
	return A5(
		_elm_lang$core$Json_Decode$object4,
		cons,
		A2(_elm_lang$core$Json_Decode_ops[':='], 'origin', _concourse$atc$Concourse_BuildEvents$decodeOrigin),
		A2(_elm_lang$core$Json_Decode_ops[':='], 'exit_status', _elm_lang$core$Json_Decode$int),
		function (_p2) {
			return A2(
				_elm_lang$core$Json_Decode$map,
				_elm_lang$core$Maybe$withDefault(_elm_lang$core$Dict$empty),
				_elm_lang$core$Json_Decode$maybe(_p2));
		}(
			A2(_elm_lang$core$Json_Decode_ops[':='], 'version', _concourse$atc$Concourse_Version$decode)),
		function (_p3) {
			return A2(
				_elm_lang$core$Json_Decode$map,
				_elm_lang$core$Maybe$withDefault(
					_elm_lang$core$Native_List.fromArray(
						[])),
				_elm_lang$core$Json_Decode$maybe(_p3));
		}(
			A2(_elm_lang$core$Json_Decode_ops[':='], 'metadata', _concourse$atc$Concourse_Metadata$decode)));
};
var _concourse$atc$Concourse_BuildEvents$BuildError = function (a) {
	return {ctor: 'BuildError', _0: a};
};
var _concourse$atc$Concourse_BuildEvents$Error = F2(
	function (a, b) {
		return {ctor: 'Error', _0: a, _1: b};
	});
var _concourse$atc$Concourse_BuildEvents$decodeErrorEvent = _elm_lang$core$Json_Decode$oneOf(
	_elm_lang$core$Native_List.fromArray(
		[
			A3(
			_elm_lang$core$Json_Decode$object2,
			_concourse$atc$Concourse_BuildEvents$Error,
			A2(_elm_lang$core$Json_Decode_ops[':='], 'origin', _concourse$atc$Concourse_BuildEvents$decodeOrigin),
			A2(_elm_lang$core$Json_Decode_ops[':='], 'message', _elm_lang$core$Json_Decode$string)),
			A2(
			_elm_lang$core$Json_Decode$object1,
			_concourse$atc$Concourse_BuildEvents$BuildError,
			A2(_elm_lang$core$Json_Decode_ops[':='], 'message', _elm_lang$core$Json_Decode$string))
		]));
var _concourse$atc$Concourse_BuildEvents$Log = F2(
	function (a, b) {
		return {ctor: 'Log', _0: a, _1: b};
	});
var _concourse$atc$Concourse_BuildEvents$FinishPut = F4(
	function (a, b, c, d) {
		return {ctor: 'FinishPut', _0: a, _1: b, _2: c, _3: d};
	});
var _concourse$atc$Concourse_BuildEvents$InitializePut = function (a) {
	return {ctor: 'InitializePut', _0: a};
};
var _concourse$atc$Concourse_BuildEvents$FinishGet = F4(
	function (a, b, c, d) {
		return {ctor: 'FinishGet', _0: a, _1: b, _2: c, _3: d};
	});
var _concourse$atc$Concourse_BuildEvents$InitializeGet = function (a) {
	return {ctor: 'InitializeGet', _0: a};
};
var _concourse$atc$Concourse_BuildEvents$FinishTask = F2(
	function (a, b) {
		return {ctor: 'FinishTask', _0: a, _1: b};
	});
var _concourse$atc$Concourse_BuildEvents$StartTask = function (a) {
	return {ctor: 'StartTask', _0: a};
};
var _concourse$atc$Concourse_BuildEvents$InitializeTask = function (a) {
	return {ctor: 'InitializeTask', _0: a};
};
var _concourse$atc$Concourse_BuildEvents$BuildStatus = F2(
	function (a, b) {
		return {ctor: 'BuildStatus', _0: a, _1: b};
	});
var _concourse$atc$Concourse_BuildEvents$decodeEvent = function (e) {
	var _p4 = e.event;
	switch (_p4) {
		case 'status':
			return A2(
				_elm_lang$core$Json_Decode$decodeValue,
				A3(
					_elm_lang$core$Json_Decode$object2,
					_concourse$atc$Concourse_BuildEvents$BuildStatus,
					A2(_elm_lang$core$Json_Decode_ops[':='], 'status', _concourse$atc$Concourse_BuildStatus$decode),
					A2(
						_elm_lang$core$Json_Decode_ops[':='],
						'time',
						A2(_elm_lang$core$Json_Decode$map, _concourse$atc$Concourse_BuildEvents$dateFromSeconds, _elm_lang$core$Json_Decode$float))),
				e.value);
		case 'log':
			return A2(
				_elm_lang$core$Json_Decode$decodeValue,
				A3(
					_elm_lang$core$Json_Decode$object2,
					_concourse$atc$Concourse_BuildEvents$Log,
					A2(_elm_lang$core$Json_Decode_ops[':='], 'origin', _concourse$atc$Concourse_BuildEvents$decodeOrigin),
					A2(_elm_lang$core$Json_Decode_ops[':='], 'payload', _elm_lang$core$Json_Decode$string)),
				e.value);
		case 'error':
			return A2(_elm_lang$core$Json_Decode$decodeValue, _concourse$atc$Concourse_BuildEvents$decodeErrorEvent, e.value);
		case 'initialize-task':
			return A2(
				_elm_lang$core$Json_Decode$decodeValue,
				A2(
					_elm_lang$core$Json_Decode$object1,
					_concourse$atc$Concourse_BuildEvents$InitializeTask,
					A2(_elm_lang$core$Json_Decode_ops[':='], 'origin', _concourse$atc$Concourse_BuildEvents$decodeOrigin)),
				e.value);
		case 'start-task':
			return A2(
				_elm_lang$core$Json_Decode$decodeValue,
				A2(
					_elm_lang$core$Json_Decode$object1,
					_concourse$atc$Concourse_BuildEvents$StartTask,
					A2(_elm_lang$core$Json_Decode_ops[':='], 'origin', _concourse$atc$Concourse_BuildEvents$decodeOrigin)),
				e.value);
		case 'finish-task':
			return A2(
				_elm_lang$core$Json_Decode$decodeValue,
				A3(
					_elm_lang$core$Json_Decode$object2,
					_concourse$atc$Concourse_BuildEvents$FinishTask,
					A2(_elm_lang$core$Json_Decode_ops[':='], 'origin', _concourse$atc$Concourse_BuildEvents$decodeOrigin),
					A2(_elm_lang$core$Json_Decode_ops[':='], 'exit_status', _elm_lang$core$Json_Decode$int)),
				e.value);
		case 'initialize-get':
			return A2(
				_elm_lang$core$Json_Decode$decodeValue,
				A2(
					_elm_lang$core$Json_Decode$object1,
					_concourse$atc$Concourse_BuildEvents$InitializeGet,
					A2(_elm_lang$core$Json_Decode_ops[':='], 'origin', _concourse$atc$Concourse_BuildEvents$decodeOrigin)),
				e.value);
		case 'finish-get':
			return A2(
				_elm_lang$core$Json_Decode$decodeValue,
				_concourse$atc$Concourse_BuildEvents$decodeFinishResource(_concourse$atc$Concourse_BuildEvents$FinishGet),
				e.value);
		case 'initialize-put':
			return A2(
				_elm_lang$core$Json_Decode$decodeValue,
				A2(
					_elm_lang$core$Json_Decode$object1,
					_concourse$atc$Concourse_BuildEvents$InitializePut,
					A2(_elm_lang$core$Json_Decode_ops[':='], 'origin', _concourse$atc$Concourse_BuildEvents$decodeOrigin)),
				e.value);
		case 'finish-put':
			return A2(
				_elm_lang$core$Json_Decode$decodeValue,
				_concourse$atc$Concourse_BuildEvents$decodeFinishResource(_concourse$atc$Concourse_BuildEvents$FinishPut),
				e.value);
		default:
			return _elm_lang$core$Result$Err(
				A2(_elm_lang$core$Basics_ops['++'], 'unknown event type: ', _p4));
	}
};
var _concourse$atc$Concourse_BuildEvents$decode = A2(_elm_lang$core$Json_Decode$customDecoder, _concourse$atc$Concourse_BuildEvents$decodeEnvelope, _concourse$atc$Concourse_BuildEvents$decodeEvent);
var _concourse$atc$Concourse_BuildEvents$parseEvent = function (data) {
	return A2(_elm_lang$core$Json_Decode$decodeString, _concourse$atc$Concourse_BuildEvents$decode, data);
};
var _concourse$atc$Concourse_BuildEvents$End = {ctor: 'End'};
var _concourse$atc$Concourse_BuildEvents$Event = function (a) {
	return {ctor: 'Event', _0: a};
};
var _concourse$atc$Concourse_BuildEvents$Errored = {ctor: 'Errored'};
var _concourse$atc$Concourse_BuildEvents$Opened = {ctor: 'Opened'};
var _concourse$atc$Concourse_BuildEvents$parseAction = function (msg) {
	var _p5 = msg;
	switch (_p5.ctor) {
		case 'Event':
			var _p8 = _p5._0.name;
			var _p7 = _p5._0.data;
			var _p6 = _p8;
			_v2_2:
			do {
				if (_p6.ctor === 'Just') {
					switch (_p6._0) {
						case 'end':
							return _concourse$atc$Concourse_BuildEvents$End;
						case 'event':
							return _concourse$atc$Concourse_BuildEvents$Event(
								_concourse$atc$Concourse_BuildEvents$parseEvent(_p7));
						default:
							break _v2_2;
					}
				} else {
					break _v2_2;
				}
			} while(false);
			return _concourse$atc$Concourse_BuildEvents$Event(
				_elm_lang$core$Result$Err(
					A2(
						_elm_lang$core$Basics_ops['++'],
						'unknown event type: ',
						A2(
							_elm_lang$core$Basics_ops['++'],
							_elm_lang$core$Basics$toString(_p8),
							A2(
								_elm_lang$core$Basics_ops['++'],
								' (data: ',
								A2(
									_elm_lang$core$Basics_ops['++'],
									_elm_lang$core$Basics$toString(_p7),
									')'))))));
		case 'Opened':
			return _concourse$atc$Concourse_BuildEvents$Opened;
		default:
			return _concourse$atc$Concourse_BuildEvents$Errored;
	}
};
var _concourse$atc$Concourse_BuildEvents$subscribe = function (build) {
	return A2(
		_concourse$atc$EventSource$listen,
		{
			ctor: '_Tuple2',
			_0: A2(
				_elm_lang$core$Basics_ops['++'],
				'/api/v1/builds/',
				A2(
					_elm_lang$core$Basics_ops['++'],
					_elm_lang$core$Basics$toString(build),
					'/events')),
			_1: _elm_lang$core$Native_List.fromArray(
				['end', 'event'])
		},
		_concourse$atc$Concourse_BuildEvents$parseAction);
};

var _concourse$atc$Concourse_BuildPlan$lazy = function (thunk) {
	return A2(
		_elm_lang$core$Json_Decode$customDecoder,
		_elm_lang$core$Json_Decode$value,
		function (js) {
			return A2(
				_elm_lang$core$Json_Decode$decodeValue,
				thunk(
					{ctor: '_Tuple0'}),
				js);
		});
};
var _concourse$atc$Concourse_BuildPlan$BuildPlan = F2(
	function (a, b) {
		return {id: a, step: b};
	});
var _concourse$atc$Concourse_BuildPlan$HookedPlan = F2(
	function (a, b) {
		return {step: a, hook: b};
	});
var _concourse$atc$Concourse_BuildPlan$Timeout = function (a) {
	return {ctor: 'Timeout', _0: a};
};
var _concourse$atc$Concourse_BuildPlan$Retry = function (a) {
	return {ctor: 'Retry', _0: a};
};
var _concourse$atc$Concourse_BuildPlan$Try = function (a) {
	return {ctor: 'Try', _0: a};
};
var _concourse$atc$Concourse_BuildPlan$Ensure = function (a) {
	return {ctor: 'Ensure', _0: a};
};
var _concourse$atc$Concourse_BuildPlan$OnFailure = function (a) {
	return {ctor: 'OnFailure', _0: a};
};
var _concourse$atc$Concourse_BuildPlan$OnSuccess = function (a) {
	return {ctor: 'OnSuccess', _0: a};
};
var _concourse$atc$Concourse_BuildPlan$Do = function (a) {
	return {ctor: 'Do', _0: a};
};
var _concourse$atc$Concourse_BuildPlan$Aggregate = function (a) {
	return {ctor: 'Aggregate', _0: a};
};
var _concourse$atc$Concourse_BuildPlan$DependentGet = function (a) {
	return {ctor: 'DependentGet', _0: a};
};
var _concourse$atc$Concourse_BuildPlan$decodeDependentGet = A2(
	_elm_lang$core$Json_Decode$object1,
	_concourse$atc$Concourse_BuildPlan$DependentGet,
	A2(_elm_lang$core$Json_Decode_ops[':='], 'name', _elm_lang$core$Json_Decode$string));
var _concourse$atc$Concourse_BuildPlan$Put = function (a) {
	return {ctor: 'Put', _0: a};
};
var _concourse$atc$Concourse_BuildPlan$decodePut = A2(
	_elm_lang$core$Json_Decode$object1,
	_concourse$atc$Concourse_BuildPlan$Put,
	A2(_elm_lang$core$Json_Decode_ops[':='], 'name', _elm_lang$core$Json_Decode$string));
var _concourse$atc$Concourse_BuildPlan$Get = F2(
	function (a, b) {
		return {ctor: 'Get', _0: a, _1: b};
	});
var _concourse$atc$Concourse_BuildPlan$decodeGet = A3(
	_elm_lang$core$Json_Decode$object2,
	_concourse$atc$Concourse_BuildPlan$Get,
	A2(_elm_lang$core$Json_Decode_ops[':='], 'name', _elm_lang$core$Json_Decode$string),
	_elm_lang$core$Json_Decode$maybe(
		A2(_elm_lang$core$Json_Decode_ops[':='], 'version', _concourse$atc$Concourse_Version$decode)));
var _concourse$atc$Concourse_BuildPlan$Task = function (a) {
	return {ctor: 'Task', _0: a};
};
var _concourse$atc$Concourse_BuildPlan$decodeTask = A2(
	_elm_lang$core$Json_Decode$object1,
	_concourse$atc$Concourse_BuildPlan$Task,
	A2(_elm_lang$core$Json_Decode_ops[':='], 'name', _elm_lang$core$Json_Decode$string));
var _concourse$atc$Concourse_BuildPlan$decodePlan = A3(
	_elm_lang$core$Json_Decode$object2,
	_concourse$atc$Concourse_BuildPlan$BuildPlan,
	A2(_elm_lang$core$Json_Decode_ops[':='], 'id', _elm_lang$core$Json_Decode$string),
	_elm_lang$core$Json_Decode$oneOf(
		_elm_lang$core$Native_List.fromArray(
			[
				A2(
				_elm_lang$core$Json_Decode_ops[':='],
				'task',
				_concourse$atc$Concourse_BuildPlan$lazy(
					function (_p0) {
						return _concourse$atc$Concourse_BuildPlan$decodeTask;
					})),
				A2(
				_elm_lang$core$Json_Decode_ops[':='],
				'get',
				_concourse$atc$Concourse_BuildPlan$lazy(
					function (_p1) {
						return _concourse$atc$Concourse_BuildPlan$decodeGet;
					})),
				A2(
				_elm_lang$core$Json_Decode_ops[':='],
				'put',
				_concourse$atc$Concourse_BuildPlan$lazy(
					function (_p2) {
						return _concourse$atc$Concourse_BuildPlan$decodePut;
					})),
				A2(
				_elm_lang$core$Json_Decode_ops[':='],
				'dependent_get',
				_concourse$atc$Concourse_BuildPlan$lazy(
					function (_p3) {
						return _concourse$atc$Concourse_BuildPlan$decodeDependentGet;
					})),
				A2(
				_elm_lang$core$Json_Decode_ops[':='],
				'aggregate',
				_concourse$atc$Concourse_BuildPlan$lazy(
					function (_p4) {
						return _concourse$atc$Concourse_BuildPlan$decodeAggregate;
					})),
				A2(
				_elm_lang$core$Json_Decode_ops[':='],
				'do',
				_concourse$atc$Concourse_BuildPlan$lazy(
					function (_p5) {
						return _concourse$atc$Concourse_BuildPlan$decodeDo;
					})),
				A2(
				_elm_lang$core$Json_Decode_ops[':='],
				'on_success',
				_concourse$atc$Concourse_BuildPlan$lazy(
					function (_p6) {
						return _concourse$atc$Concourse_BuildPlan$decodeOnSuccess;
					})),
				A2(
				_elm_lang$core$Json_Decode_ops[':='],
				'on_failure',
				_concourse$atc$Concourse_BuildPlan$lazy(
					function (_p7) {
						return _concourse$atc$Concourse_BuildPlan$decodeOnFailure;
					})),
				A2(
				_elm_lang$core$Json_Decode_ops[':='],
				'ensure',
				_concourse$atc$Concourse_BuildPlan$lazy(
					function (_p8) {
						return _concourse$atc$Concourse_BuildPlan$decodeEnsure;
					})),
				A2(
				_elm_lang$core$Json_Decode_ops[':='],
				'try',
				_concourse$atc$Concourse_BuildPlan$lazy(
					function (_p9) {
						return _concourse$atc$Concourse_BuildPlan$decodeTry;
					})),
				A2(
				_elm_lang$core$Json_Decode_ops[':='],
				'retry',
				_concourse$atc$Concourse_BuildPlan$lazy(
					function (_p10) {
						return _concourse$atc$Concourse_BuildPlan$decodeRetry;
					})),
				A2(
				_elm_lang$core$Json_Decode_ops[':='],
				'timeout',
				_concourse$atc$Concourse_BuildPlan$lazy(
					function (_p11) {
						return _concourse$atc$Concourse_BuildPlan$decodeTimeout;
					}))
			])));
var _concourse$atc$Concourse_BuildPlan$decodeAggregate = A2(
	_elm_lang$core$Json_Decode$object1,
	_concourse$atc$Concourse_BuildPlan$Aggregate,
	_elm_lang$core$Json_Decode$array(
		_concourse$atc$Concourse_BuildPlan$lazy(
			function (_p12) {
				return _concourse$atc$Concourse_BuildPlan$decodePlan;
			})));
var _concourse$atc$Concourse_BuildPlan$decodeDo = A2(
	_elm_lang$core$Json_Decode$object1,
	_concourse$atc$Concourse_BuildPlan$Do,
	_elm_lang$core$Json_Decode$array(
		_concourse$atc$Concourse_BuildPlan$lazy(
			function (_p13) {
				return _concourse$atc$Concourse_BuildPlan$decodePlan;
			})));
var _concourse$atc$Concourse_BuildPlan$decodeEnsure = A2(
	_elm_lang$core$Json_Decode$map,
	_concourse$atc$Concourse_BuildPlan$Ensure,
	A3(
		_elm_lang$core$Json_Decode$object2,
		_concourse$atc$Concourse_BuildPlan$HookedPlan,
		A2(
			_elm_lang$core$Json_Decode_ops[':='],
			'step',
			_concourse$atc$Concourse_BuildPlan$lazy(
				function (_p14) {
					return _concourse$atc$Concourse_BuildPlan$decodePlan;
				})),
		A2(
			_elm_lang$core$Json_Decode_ops[':='],
			'ensure',
			_concourse$atc$Concourse_BuildPlan$lazy(
				function (_p15) {
					return _concourse$atc$Concourse_BuildPlan$decodePlan;
				}))));
var _concourse$atc$Concourse_BuildPlan$decodeOnFailure = A2(
	_elm_lang$core$Json_Decode$map,
	_concourse$atc$Concourse_BuildPlan$OnFailure,
	A3(
		_elm_lang$core$Json_Decode$object2,
		_concourse$atc$Concourse_BuildPlan$HookedPlan,
		A2(
			_elm_lang$core$Json_Decode_ops[':='],
			'step',
			_concourse$atc$Concourse_BuildPlan$lazy(
				function (_p16) {
					return _concourse$atc$Concourse_BuildPlan$decodePlan;
				})),
		A2(
			_elm_lang$core$Json_Decode_ops[':='],
			'on_failure',
			_concourse$atc$Concourse_BuildPlan$lazy(
				function (_p17) {
					return _concourse$atc$Concourse_BuildPlan$decodePlan;
				}))));
var _concourse$atc$Concourse_BuildPlan$decodeOnSuccess = A2(
	_elm_lang$core$Json_Decode$map,
	_concourse$atc$Concourse_BuildPlan$OnSuccess,
	A3(
		_elm_lang$core$Json_Decode$object2,
		_concourse$atc$Concourse_BuildPlan$HookedPlan,
		A2(
			_elm_lang$core$Json_Decode_ops[':='],
			'step',
			_concourse$atc$Concourse_BuildPlan$lazy(
				function (_p18) {
					return _concourse$atc$Concourse_BuildPlan$decodePlan;
				})),
		A2(
			_elm_lang$core$Json_Decode_ops[':='],
			'on_success',
			_concourse$atc$Concourse_BuildPlan$lazy(
				function (_p19) {
					return _concourse$atc$Concourse_BuildPlan$decodePlan;
				}))));
var _concourse$atc$Concourse_BuildPlan$decodeRetry = A2(
	_elm_lang$core$Json_Decode$object1,
	_concourse$atc$Concourse_BuildPlan$Retry,
	_elm_lang$core$Json_Decode$array(
		_concourse$atc$Concourse_BuildPlan$lazy(
			function (_p20) {
				return _concourse$atc$Concourse_BuildPlan$decodePlan;
			})));
var _concourse$atc$Concourse_BuildPlan$decodeTimeout = A2(
	_elm_lang$core$Json_Decode$object1,
	_concourse$atc$Concourse_BuildPlan$Timeout,
	A2(
		_elm_lang$core$Json_Decode_ops[':='],
		'step',
		_concourse$atc$Concourse_BuildPlan$lazy(
			function (_p21) {
				return _concourse$atc$Concourse_BuildPlan$decodePlan;
			})));
var _concourse$atc$Concourse_BuildPlan$decodeTry = A2(
	_elm_lang$core$Json_Decode$object1,
	_concourse$atc$Concourse_BuildPlan$Try,
	A2(
		_elm_lang$core$Json_Decode_ops[':='],
		'step',
		_concourse$atc$Concourse_BuildPlan$lazy(
			function (_p22) {
				return _concourse$atc$Concourse_BuildPlan$decodePlan;
			})));
var _concourse$atc$Concourse_BuildPlan$decode = A2(
	_elm_lang$core$Json_Decode$at,
	_elm_lang$core$Native_List.fromArray(
		['plan']),
	_concourse$atc$Concourse_BuildPlan$decodePlan);
var _concourse$atc$Concourse_BuildPlan$fetch = function (buildId) {
	return A2(
		_evancz$elm_http$Http$get,
		_concourse$atc$Concourse_BuildPlan$decode,
		A2(
			_elm_lang$core$Basics_ops['++'],
			'/api/v1/builds/',
			A2(
				_elm_lang$core$Basics_ops['++'],
				_elm_lang$core$Basics$toString(buildId),
				'/plan')));
};

var _concourse$atc$Concourse_BuildResources$empty = {
	inputs: _elm_lang$core$Native_List.fromArray(
		[]),
	outputs: _elm_lang$core$Native_List.fromArray(
		[])
};
var _concourse$atc$Concourse_BuildResources$BuildResources = F2(
	function (a, b) {
		return {inputs: a, outputs: b};
	});
var _concourse$atc$Concourse_BuildResources$BuildInput = F6(
	function (a, b, c, d, e, f) {
		return {name: a, resource: b, type$: c, version: d, metadata: e, firstOccurrence: f};
	});
var _concourse$atc$Concourse_BuildResources$decodeInput = A7(
	_elm_lang$core$Json_Decode$object6,
	_concourse$atc$Concourse_BuildResources$BuildInput,
	A2(_elm_lang$core$Json_Decode_ops[':='], 'name', _elm_lang$core$Json_Decode$string),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'resource', _elm_lang$core$Json_Decode$string),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'type', _elm_lang$core$Json_Decode$string),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'version', _concourse$atc$Concourse_Version$decode),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'metadata', _concourse$atc$Concourse_Metadata$decode),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'first_occurrence', _elm_lang$core$Json_Decode$bool));
var _concourse$atc$Concourse_BuildResources$BuildOutput = F2(
	function (a, b) {
		return {resource: a, version: b};
	});
var _concourse$atc$Concourse_BuildResources$decodeOutput = A3(
	_elm_lang$core$Json_Decode$object2,
	_concourse$atc$Concourse_BuildResources$BuildOutput,
	A2(_elm_lang$core$Json_Decode_ops[':='], 'resource', _elm_lang$core$Json_Decode$string),
	A2(
		_elm_lang$core$Json_Decode_ops[':='],
		'version',
		_elm_lang$core$Json_Decode$dict(_elm_lang$core$Json_Decode$string)));
var _concourse$atc$Concourse_BuildResources$decode = A3(
	_elm_lang$core$Json_Decode$object2,
	_concourse$atc$Concourse_BuildResources$BuildResources,
	A2(
		_elm_lang$core$Json_Decode_ops[':='],
		'inputs',
		_elm_lang$core$Json_Decode$list(_concourse$atc$Concourse_BuildResources$decodeInput)),
	A2(
		_elm_lang$core$Json_Decode_ops[':='],
		'outputs',
		_elm_lang$core$Json_Decode$list(_concourse$atc$Concourse_BuildResources$decodeOutput)));
var _concourse$atc$Concourse_BuildResources$fetch = function (buildId) {
	return A2(
		_evancz$elm_http$Http$get,
		_concourse$atc$Concourse_BuildResources$decode,
		A2(
			_elm_lang$core$Basics_ops['++'],
			'/api/v1/builds/',
			A2(
				_elm_lang$core$Basics_ops['++'],
				_elm_lang$core$Basics$toString(buildId),
				'/resources')));
};

var _concourse$atc$LoadingIndicator$view = A2(
	_elm_lang$html$Html$div,
	_elm_lang$core$Native_List.fromArray(
		[
			_elm_lang$html$Html_Attributes$class('build-step')
		]),
	_elm_lang$core$Native_List.fromArray(
		[
			A2(
			_elm_lang$html$Html$div,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$class('header')
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(
					_elm_lang$html$Html$i,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('left fa fa-fw fa-spin fa-circle-o-notch')
						]),
					_elm_lang$core$Native_List.fromArray(
						[])),
					A2(
					_elm_lang$html$Html$h3,
					_elm_lang$core$Native_List.fromArray(
						[]),
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html$text('loading')
						]))
				]))
		]));

var _concourse$atc$Focus$update = F3(
	function (_p0, f, big) {
		var _p1 = _p0;
		return A2(_p1._0.update, f, big);
	});
var _concourse$atc$Focus$set = F3(
	function (_p2, small, big) {
		var _p3 = _p2;
		return A2(
			_p3._0.update,
			_elm_lang$core$Basics$always(small),
			big);
	});
var _concourse$atc$Focus$get = F2(
	function (_p4, big) {
		var _p5 = _p4;
		return _p5._0.get(big);
	});
var _concourse$atc$Focus$Focus = function (a) {
	return {ctor: 'Focus', _0: a};
};
var _concourse$atc$Focus$create = F2(
	function (get, update) {
		return _concourse$atc$Focus$Focus(
			{get: get, update: update});
	});
var _concourse$atc$Focus_ops = _concourse$atc$Focus_ops || {};
_concourse$atc$Focus_ops['=>'] = F2(
	function (_p7, _p6) {
		var _p8 = _p7;
		var _p11 = _p8._0;
		var _p9 = _p6;
		var _p10 = _p9._0;
		var update = F2(
			function (f, big) {
				return A2(
					_p11.update,
					_p10.update(f),
					big);
			});
		var get = function (big) {
			return _p10.get(
				_p11.get(big));
		};
		return _concourse$atc$Focus$Focus(
			{get: get, update: update});
	});

var _concourse$atc$DictView$viewPair = function (_p0) {
	var _p1 = _p0;
	return A2(
		_elm_lang$html$Html$tr,
		_elm_lang$core$Native_List.fromArray(
			[]),
		_elm_lang$core$Native_List.fromArray(
			[
				A2(
				_elm_lang$html$Html$td,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('dict-key')
					]),
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html$text(_p1._0)
					])),
				A2(
				_elm_lang$html$Html$td,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('dict-value')
					]),
				_elm_lang$core$Native_List.fromArray(
					[_p1._1]))
			]));
};
var _concourse$atc$DictView$view = function (dict) {
	return A2(
		_elm_lang$html$Html$table,
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html_Attributes$class('dictionary')
			]),
		A2(
			_elm_lang$core$List$map,
			_concourse$atc$DictView$viewPair,
			_elm_lang$core$Dict$toList(dict)));
};

var _concourse$atc$StepTree$viewStepState = F2(
	function (state, finished) {
		var _p0 = state;
		switch (_p0.ctor) {
			case 'StepStatePending':
				var icon = finished ? 'fa-circle' : 'fa-circle-o-notch';
				return A2(
					_elm_lang$html$Html$i,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class(
							A2(_elm_lang$core$Basics_ops['++'], 'right fa fa-fw ', icon))
						]),
					_elm_lang$core$Native_List.fromArray(
						[]));
			case 'StepStateRunning':
				return A2(
					_elm_lang$html$Html$i,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('right fa fa-fw fa-spin fa-circle-o-notch')
						]),
					_elm_lang$core$Native_List.fromArray(
						[]));
			case 'StepStateSucceeded':
				return A2(
					_elm_lang$html$Html$i,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('right succeeded fa fa-fw fa-check')
						]),
					_elm_lang$core$Native_List.fromArray(
						[]));
			case 'StepStateFailed':
				return A2(
					_elm_lang$html$Html$i,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('right failed fa fa-fw fa-times')
						]),
					_elm_lang$core$Native_List.fromArray(
						[]));
			default:
				return A2(
					_elm_lang$html$Html$i,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('right errored fa fa-fw fa-exclamation-triangle')
						]),
					_elm_lang$core$Native_List.fromArray(
						[]));
		}
	});
var _concourse$atc$StepTree$typeIcon = function (fa) {
	return A2(
		_elm_lang$html$Html$i,
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html_Attributes$class(
				A2(_elm_lang$core$Basics_ops['++'], 'left fa fa-fw ', fa))
			]),
		_elm_lang$core$Native_List.fromArray(
			[]));
};
var _concourse$atc$StepTree$viewMetadata = function (metadata) {
	return function (_p1) {
		return _concourse$atc$DictView$view(
			_elm_lang$core$Dict$fromList(_p1));
	}(
		A2(
			_elm_lang$core$List$map,
			function (_p2) {
				var _p3 = _p2;
				return {
					ctor: '_Tuple2',
					_0: _p3.name,
					_1: A2(
						_elm_lang$html$Html$pre,
						_elm_lang$core$Native_List.fromArray(
							[]),
						_elm_lang$core$Native_List.fromArray(
							[
								_elm_lang$html$Html$text(_p3.value)
							]))
				};
			},
			metadata));
};
var _concourse$atc$StepTree$viewVersion = function (version) {
	return function (_p4) {
		return _concourse$atc$DictView$view(
			A2(
				_elm_lang$core$Dict$map,
				F2(
					function (_p5, s) {
						return _elm_lang$html$Html$text(s);
					}),
				_p4));
	}(
		A2(_elm_lang$core$Maybe$withDefault, _elm_lang$core$Dict$empty, version));
};
var _concourse$atc$StepTree$getMultiStepIndex = F2(
	function (idx, tree) {
		var steps = function () {
			var _p6 = tree;
			switch (_p6.ctor) {
				case 'Aggregate':
					return _p6._0;
				case 'Do':
					return _p6._0;
				case 'Retry':
					return _p6._1;
				default:
					return _elm_lang$core$Native_Utils.crashCase(
						'StepTree',
						{
							start: {line: 404, column: 7},
							end: {line: 415, column: 35}
						},
						_p6)('impossible');
			}
		}();
		var _p8 = A2(_elm_lang$core$Array$get, idx, steps);
		if (_p8.ctor === 'Just') {
			return _p8._0;
		} else {
			return _elm_lang$core$Native_Utils.crashCase(
				'StepTree',
				{
					start: {line: 417, column: 5},
					end: {line: 422, column: 33}
				},
				_p8)('impossible');
		}
	});
var _concourse$atc$StepTree$getHook = function (tree) {
	var _p10 = tree;
	switch (_p10.ctor) {
		case 'OnSuccess':
			return _p10._0.hook;
		case 'OnFailure':
			return _p10._0.hook;
		case 'Ensure':
			return _p10._0.hook;
		default:
			return _elm_lang$core$Native_Utils.crashCase(
				'StepTree',
				{
					start: {line: 372, column: 3},
					end: {line: 383, column: 31}
				},
				_p10)('impossible');
	}
};
var _concourse$atc$StepTree$getStep = function (tree) {
	var _p12 = tree;
	switch (_p12.ctor) {
		case 'OnSuccess':
			return _p12._0.step;
		case 'OnFailure':
			return _p12._0.step;
		case 'Ensure':
			return _p12._0.step;
		case 'Try':
			return _p12._0;
		case 'Timeout':
			return _p12._0;
		default:
			return _elm_lang$core$Native_Utils.crashCase(
				'StepTree',
				{
					start: {line: 326, column: 3},
					end: {line: 343, column: 31}
				},
				_p12)('impossible');
	}
};
var _concourse$atc$StepTree$updateAt = F3(
	function (id, update, root) {
		var _p14 = A2(_elm_lang$core$Dict$get, id, root.foci);
		if (_p14.ctor === 'Nothing') {
			return _elm_lang$core$Native_Utils.crashCase(
				'StepTree',
				{
					start: {line: 248, column: 3},
					end: {line: 253, column: 60}
				},
				_p14)(
				A2(
					_elm_lang$core$Basics_ops['++'],
					'updateAt: id ',
					A2(_elm_lang$core$Basics_ops['++'], id, ' not found')));
		} else {
			return _elm_lang$core$Native_Utils.update(
				root,
				{
					tree: A3(_concourse$atc$Focus$update, _p14._0, update, root.tree)
				});
		}
	});
var _concourse$atc$StepTree$isFirstOccurrence = F2(
	function (resources, step) {
		isFirstOccurrence:
		while (true) {
			var _p16 = resources;
			if (_p16.ctor === '[]') {
				return false;
			} else {
				if (_elm_lang$core$Native_Utils.eq(_p16._0.name, step)) {
					return _p16._0.firstOccurrence;
				} else {
					var _v8 = _p16._1,
						_v9 = step;
					resources = _v8;
					step = _v9;
					continue isFirstOccurrence;
				}
			}
		}
	});
var _concourse$atc$StepTree$setupGetStep = F4(
	function (resources, name, version, step) {
		return _elm_lang$core$Native_Utils.update(
			step,
			{
				version: version,
				firstOccurrence: A2(_concourse$atc$StepTree$isFirstOccurrence, resources.inputs, name)
			});
	});
var _concourse$atc$StepTree$HookedStep = F2(
	function (a, b) {
		return {step: a, hook: b};
	});
var _concourse$atc$StepTree$Step = F9(
	function (a, b, c, d, e, f, g, h, i) {
		return {id: a, name: b, state: c, log: d, error: e, expanded: f, version: g, metadata: h, firstOccurrence: i};
	});
var _concourse$atc$StepTree$Model = F3(
	function (a, b, c) {
		return {tree: a, foci: b, finished: c};
	});
var _concourse$atc$StepTree$MetadataField = F2(
	function (a, b) {
		return {name: a, value: b};
	});
var _concourse$atc$StepTree$Timeout = function (a) {
	return {ctor: 'Timeout', _0: a};
};
var _concourse$atc$StepTree$Retry = F4(
	function (a, b, c, d) {
		return {ctor: 'Retry', _0: a, _1: b, _2: c, _3: d};
	});
var _concourse$atc$StepTree$Try = function (a) {
	return {ctor: 'Try', _0: a};
};
var _concourse$atc$StepTree$Ensure = function (a) {
	return {ctor: 'Ensure', _0: a};
};
var _concourse$atc$StepTree$OnFailure = function (a) {
	return {ctor: 'OnFailure', _0: a};
};
var _concourse$atc$StepTree$OnSuccess = function (a) {
	return {ctor: 'OnSuccess', _0: a};
};
var _concourse$atc$StepTree$updateStep = F2(
	function (update, tree) {
		var _p17 = tree;
		switch (_p17.ctor) {
			case 'OnSuccess':
				var _p18 = _p17._0;
				return _concourse$atc$StepTree$OnSuccess(
					_elm_lang$core$Native_Utils.update(
						_p18,
						{
							step: update(_p18.step)
						}));
			case 'OnFailure':
				var _p19 = _p17._0;
				return _concourse$atc$StepTree$OnFailure(
					_elm_lang$core$Native_Utils.update(
						_p19,
						{
							step: update(_p19.step)
						}));
			case 'Ensure':
				var _p20 = _p17._0;
				return _concourse$atc$StepTree$Ensure(
					_elm_lang$core$Native_Utils.update(
						_p20,
						{
							step: update(_p20.step)
						}));
			case 'Try':
				return _concourse$atc$StepTree$Try(
					update(_p17._0));
			case 'Timeout':
				return _concourse$atc$StepTree$Timeout(
					update(_p17._0));
			default:
				return _elm_lang$core$Native_Utils.crashCase(
					'StepTree',
					{
						start: {line: 347, column: 3},
						end: {line: 364, column: 31}
					},
					_p17)('impossible');
		}
	});
var _concourse$atc$StepTree$wrapStep = F2(
	function (id, subFocus) {
		return A2(
			_concourse$atc$Focus_ops['=>'],
			A2(_concourse$atc$Focus$create, _concourse$atc$StepTree$getStep, _concourse$atc$StepTree$updateStep),
			subFocus);
	});
var _concourse$atc$StepTree$updateHook = F2(
	function (update, tree) {
		var _p22 = tree;
		switch (_p22.ctor) {
			case 'OnSuccess':
				var _p23 = _p22._0;
				return _concourse$atc$StepTree$OnSuccess(
					_elm_lang$core$Native_Utils.update(
						_p23,
						{
							hook: update(_p23.hook)
						}));
			case 'OnFailure':
				var _p24 = _p22._0;
				return _concourse$atc$StepTree$OnFailure(
					_elm_lang$core$Native_Utils.update(
						_p24,
						{
							hook: update(_p24.hook)
						}));
			case 'Ensure':
				var _p25 = _p22._0;
				return _concourse$atc$StepTree$Ensure(
					_elm_lang$core$Native_Utils.update(
						_p25,
						{
							hook: update(_p25.hook)
						}));
			default:
				return _elm_lang$core$Native_Utils.crashCase(
					'StepTree',
					{
						start: {line: 387, column: 3},
						end: {line: 398, column: 31}
					},
					_p22)('impossible');
		}
	});
var _concourse$atc$StepTree$wrapHook = F2(
	function (id, subFocus) {
		return A2(
			_concourse$atc$Focus_ops['=>'],
			A2(_concourse$atc$Focus$create, _concourse$atc$StepTree$getHook, _concourse$atc$StepTree$updateHook),
			subFocus);
	});
var _concourse$atc$StepTree$Do = function (a) {
	return {ctor: 'Do', _0: a};
};
var _concourse$atc$StepTree$Aggregate = function (a) {
	return {ctor: 'Aggregate', _0: a};
};
var _concourse$atc$StepTree$DependentGet = function (a) {
	return {ctor: 'DependentGet', _0: a};
};
var _concourse$atc$StepTree$Put = function (a) {
	return {ctor: 'Put', _0: a};
};
var _concourse$atc$StepTree$Get = function (a) {
	return {ctor: 'Get', _0: a};
};
var _concourse$atc$StepTree$Task = function (a) {
	return {ctor: 'Task', _0: a};
};
var _concourse$atc$StepTree$map = F2(
	function (f, tree) {
		var _p27 = tree;
		switch (_p27.ctor) {
			case 'Task':
				return _concourse$atc$StepTree$Task(
					f(_p27._0));
			case 'Get':
				return _concourse$atc$StepTree$Get(
					f(_p27._0));
			case 'Put':
				return _concourse$atc$StepTree$Put(
					f(_p27._0));
			case 'DependentGet':
				return _concourse$atc$StepTree$DependentGet(
					f(_p27._0));
			default:
				return tree;
		}
	});
var _concourse$atc$StepTree$User = {ctor: 'User'};
var _concourse$atc$StepTree$focusRetry = F2(
	function (tab, tree) {
		var _p28 = tree;
		if (_p28.ctor === 'Retry') {
			return A4(_concourse$atc$StepTree$Retry, _p28._0, _p28._1, tab, _concourse$atc$StepTree$User);
		} else {
			return _elm_lang$core$Native_Utils.crashCase(
				'StepTree',
				{
					start: {line: 239, column: 3},
					end: {line: 244, column: 53}
				},
				_p28)('impossible (non-retry tab focus)');
		}
	});
var _concourse$atc$StepTree$Auto = {ctor: 'Auto'};
var _concourse$atc$StepTree$setMultiStepIndex = F3(
	function (idx, update, tree) {
		var _p30 = tree;
		switch (_p30.ctor) {
			case 'Aggregate':
				return _concourse$atc$StepTree$Aggregate(
					A3(
						_elm_lang$core$Array$set,
						idx,
						update(
							A2(_concourse$atc$StepTree$getMultiStepIndex, idx, tree)),
						_p30._0));
			case 'Do':
				return _concourse$atc$StepTree$Do(
					A3(
						_elm_lang$core$Array$set,
						idx,
						update(
							A2(_concourse$atc$StepTree$getMultiStepIndex, idx, tree)),
						_p30._0));
			case 'Retry':
				var _p32 = _p30._0;
				var updatedSteps = A3(
					_elm_lang$core$Array$set,
					idx,
					update(
						A2(_concourse$atc$StepTree$getMultiStepIndex, idx, tree)),
					_p30._1);
				var _p31 = _p30._3;
				if (_p31.ctor === 'Auto') {
					return A4(_concourse$atc$StepTree$Retry, _p32, updatedSteps, idx + 1, _concourse$atc$StepTree$Auto);
				} else {
					return A4(_concourse$atc$StepTree$Retry, _p32, updatedSteps, _p30._2, _concourse$atc$StepTree$User);
				}
			default:
				return _elm_lang$core$Native_Utils.crashCase(
					'StepTree',
					{
						start: {line: 426, column: 3},
						end: {line: 444, column: 31}
					},
					_p30)('impossible');
		}
	});
var _concourse$atc$StepTree$wrapMultiStep = function (i) {
	return _elm_lang$core$Dict$map(
		F2(
			function (_p34, focus) {
				return A2(
					_concourse$atc$Focus_ops['=>'],
					A2(
						_concourse$atc$Focus$create,
						_concourse$atc$StepTree$getMultiStepIndex(i),
						_concourse$atc$StepTree$setMultiStepIndex(i)),
					focus);
			}));
};
var _concourse$atc$StepTree$SwitchTab = F2(
	function (a, b) {
		return {ctor: 'SwitchTab', _0: a, _1: b};
	});
var _concourse$atc$StepTree$Finished = {ctor: 'Finished'};
var _concourse$atc$StepTree$ToggleStep = function (a) {
	return {ctor: 'ToggleStep', _0: a};
};
var _concourse$atc$StepTree$StepStateErrored = {ctor: 'StepStateErrored'};
var _concourse$atc$StepTree$StepStateFailed = {ctor: 'StepStateFailed'};
var _concourse$atc$StepTree$StepStateSucceeded = {ctor: 'StepStateSucceeded'};
var _concourse$atc$StepTree$StepStateRunning = {ctor: 'StepStateRunning'};
var _concourse$atc$StepTree$StepStatePending = {ctor: 'StepStatePending'};
var _concourse$atc$StepTree$initBottom = F3(
	function (create, id, name) {
		var step = {
			id: id,
			name: name,
			state: _concourse$atc$StepTree$StepStatePending,
			log: _vito$elm_ansi$Ansi_Log$init(_vito$elm_ansi$Ansi_Log$Cooked),
			error: _elm_lang$core$Maybe$Nothing,
			expanded: _elm_lang$core$Maybe$Nothing,
			version: _elm_lang$core$Maybe$Nothing,
			metadata: _elm_lang$core$Native_List.fromArray(
				[]),
			firstOccurrence: false
		};
		return {
			tree: create(step),
			foci: A2(
				_elm_lang$core$Dict$singleton,
				id,
				A2(_concourse$atc$Focus$create, _elm_lang$core$Basics$identity, _elm_lang$core$Basics$identity)),
			finished: false
		};
	});
var _concourse$atc$StepTree$init = F2(
	function (resources, plan) {
		var _p35 = plan.step;
		switch (_p35.ctor) {
			case 'Task':
				return A3(_concourse$atc$StepTree$initBottom, _concourse$atc$StepTree$Task, plan.id, _p35._0);
			case 'Get':
				var _p37 = _p35._0;
				return A3(
					_concourse$atc$StepTree$initBottom,
					function (_p36) {
						return _concourse$atc$StepTree$Get(
							A4(_concourse$atc$StepTree$setupGetStep, resources, _p37, _p35._1, _p36));
					},
					plan.id,
					_p37);
			case 'Put':
				return A3(_concourse$atc$StepTree$initBottom, _concourse$atc$StepTree$Put, plan.id, _p35._0);
			case 'DependentGet':
				return A3(_concourse$atc$StepTree$initBottom, _concourse$atc$StepTree$DependentGet, plan.id, _p35._0);
			case 'Aggregate':
				var inited = A2(
					_elm_lang$core$Array$map,
					_concourse$atc$StepTree$init(resources),
					_p35._0);
				var trees = A2(
					_elm_lang$core$Array$map,
					function (_) {
						return _.tree;
					},
					inited);
				var subFoci = A2(
					_elm_lang$core$Array$map,
					function (_) {
						return _.foci;
					},
					inited);
				var wrappedSubFoci = A2(_elm_lang$core$Array$indexedMap, _concourse$atc$StepTree$wrapMultiStep, subFoci);
				var foci = A3(_elm_lang$core$Array$foldr, _elm_lang$core$Dict$union, _elm_lang$core$Dict$empty, wrappedSubFoci);
				return A3(
					_concourse$atc$StepTree$Model,
					_concourse$atc$StepTree$Aggregate(trees),
					foci,
					false);
			case 'Do':
				var inited = A2(
					_elm_lang$core$Array$map,
					_concourse$atc$StepTree$init(resources),
					_p35._0);
				var trees = A2(
					_elm_lang$core$Array$map,
					function (_) {
						return _.tree;
					},
					inited);
				var subFoci = A2(
					_elm_lang$core$Array$map,
					function (_) {
						return _.foci;
					},
					inited);
				var wrappedSubFoci = A2(_elm_lang$core$Array$indexedMap, _concourse$atc$StepTree$wrapMultiStep, subFoci);
				var foci = A3(_elm_lang$core$Array$foldr, _elm_lang$core$Dict$union, _elm_lang$core$Dict$empty, wrappedSubFoci);
				return A3(
					_concourse$atc$StepTree$Model,
					_concourse$atc$StepTree$Do(trees),
					foci,
					false);
			case 'OnSuccess':
				return A3(_concourse$atc$StepTree$initHookedStep, resources, _concourse$atc$StepTree$OnSuccess, _p35._0);
			case 'OnFailure':
				return A3(_concourse$atc$StepTree$initHookedStep, resources, _concourse$atc$StepTree$OnFailure, _p35._0);
			case 'Ensure':
				return A3(_concourse$atc$StepTree$initHookedStep, resources, _concourse$atc$StepTree$Ensure, _p35._0);
			case 'Try':
				return A3(_concourse$atc$StepTree$initWrappedStep, resources, _concourse$atc$StepTree$Try, _p35._0);
			case 'Retry':
				var selfFoci = A2(
					_elm_lang$core$Dict$singleton,
					plan.id,
					A2(_concourse$atc$Focus$create, _elm_lang$core$Basics$identity, _elm_lang$core$Basics$identity));
				var inited = A2(
					_elm_lang$core$Array$map,
					_concourse$atc$StepTree$init(resources),
					_p35._0);
				var trees = A2(
					_elm_lang$core$Array$map,
					function (_) {
						return _.tree;
					},
					inited);
				var subFoci = A2(
					_elm_lang$core$Array$map,
					function (_) {
						return _.foci;
					},
					inited);
				var wrappedSubFoci = A2(_elm_lang$core$Array$indexedMap, _concourse$atc$StepTree$wrapMultiStep, subFoci);
				var foci = A3(_elm_lang$core$Array$foldr, _elm_lang$core$Dict$union, selfFoci, wrappedSubFoci);
				return A3(
					_concourse$atc$StepTree$Model,
					A4(_concourse$atc$StepTree$Retry, plan.id, trees, 1, _concourse$atc$StepTree$Auto),
					foci,
					false);
			default:
				return A3(_concourse$atc$StepTree$initWrappedStep, resources, _concourse$atc$StepTree$Timeout, _p35._0);
		}
	});
var _concourse$atc$StepTree$initHookedStep = F3(
	function (resources, create, hookedPlan) {
		var hookModel = A2(_concourse$atc$StepTree$init, resources, hookedPlan.hook);
		var stepModel = A2(_concourse$atc$StepTree$init, resources, hookedPlan.step);
		return {
			tree: create(
				{step: stepModel.tree, hook: hookModel.tree}),
			foci: A2(
				_elm_lang$core$Dict$union,
				A2(_elm_lang$core$Dict$map, _concourse$atc$StepTree$wrapStep, stepModel.foci),
				A2(_elm_lang$core$Dict$map, _concourse$atc$StepTree$wrapHook, hookModel.foci)),
			finished: stepModel.finished
		};
	});
var _concourse$atc$StepTree$initWrappedStep = F3(
	function (resources, create, plan) {
		var _p38 = A2(_concourse$atc$StepTree$init, resources, plan);
		var tree = _p38.tree;
		var foci = _p38.foci;
		return {
			tree: create(tree),
			foci: A2(_elm_lang$core$Dict$map, _concourse$atc$StepTree$wrapStep, foci),
			finished: false
		};
	});
var _concourse$atc$StepTree$isActive = F2(
	function (x, y) {
		return !_elm_lang$core$Native_Utils.eq(x, y);
	})(_concourse$atc$StepTree$StepStatePending);
var _concourse$atc$StepTree$stepIsActive = function (_p39) {
	return _concourse$atc$StepTree$isActive(
		function (_) {
			return _.state;
		}(_p39));
};
var _concourse$atc$StepTree$treeIsActive = function (tree) {
	treeIsActive:
	while (true) {
		var _p40 = tree;
		switch (_p40.ctor) {
			case 'Aggregate':
				return A2(
					_elm_lang$core$List$any,
					_concourse$atc$StepTree$treeIsActive,
					_elm_lang$core$Array$toList(_p40._0));
			case 'Do':
				return A2(
					_elm_lang$core$List$any,
					_concourse$atc$StepTree$treeIsActive,
					_elm_lang$core$Array$toList(_p40._0));
			case 'OnSuccess':
				var _v18 = _p40._0.step;
				tree = _v18;
				continue treeIsActive;
			case 'OnFailure':
				var _v19 = _p40._0.step;
				tree = _v19;
				continue treeIsActive;
			case 'Ensure':
				var _v20 = _p40._0.step;
				tree = _v20;
				continue treeIsActive;
			case 'Try':
				var _v21 = _p40._0;
				tree = _v21;
				continue treeIsActive;
			case 'Timeout':
				var _v22 = _p40._0;
				tree = _v22;
				continue treeIsActive;
			case 'Retry':
				return A2(
					_elm_lang$core$List$any,
					_concourse$atc$StepTree$treeIsActive,
					_elm_lang$core$Array$toList(_p40._1));
			case 'Task':
				return _concourse$atc$StepTree$stepIsActive(_p40._0);
			case 'Get':
				return _concourse$atc$StepTree$stepIsActive(_p40._0);
			case 'Put':
				return _concourse$atc$StepTree$stepIsActive(_p40._0);
			default:
				return _concourse$atc$StepTree$stepIsActive(_p40._0);
		}
	}
};
var _concourse$atc$StepTree$viewTab = F4(
	function (id, currentTab, idx, step) {
		var tab = idx + 1;
		return A2(
			_elm_lang$html$Html$li,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$classList(
					_elm_lang$core$Native_List.fromArray(
						[
							{
							ctor: '_Tuple2',
							_0: 'current',
							_1: _elm_lang$core$Native_Utils.eq(currentTab, tab)
						},
							{
							ctor: '_Tuple2',
							_0: 'inactive',
							_1: _elm_lang$core$Basics$not(
								_concourse$atc$StepTree$treeIsActive(step))
						}
						]))
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(
					_elm_lang$html$Html$span,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Events$onClick(
							A2(_concourse$atc$StepTree$SwitchTab, id, tab))
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html$text(
							_elm_lang$core$Basics$toString(tab))
						]))
				]));
	});
var _concourse$atc$StepTree$autoExpanded = function (state) {
	return _concourse$atc$StepTree$isActive(state) && (!_elm_lang$core$Native_Utils.eq(state, _concourse$atc$StepTree$StepStateSucceeded));
};
var _concourse$atc$StepTree$toggleExpanded = function (_p41) {
	var _p42 = _p41;
	return _elm_lang$core$Maybe$Just(
		_elm_lang$core$Basics$not(
			A2(
				_elm_lang$core$Maybe$withDefault,
				_concourse$atc$StepTree$autoExpanded(_p42.state),
				_p42.expanded)));
};
var _concourse$atc$StepTree$update = F2(
	function (action, root) {
		var _p43 = action;
		switch (_p43.ctor) {
			case 'ToggleStep':
				return A3(
					_concourse$atc$StepTree$updateAt,
					_p43._0,
					_concourse$atc$StepTree$map(
						function (step) {
							return _elm_lang$core$Native_Utils.update(
								step,
								{
									expanded: _concourse$atc$StepTree$toggleExpanded(step)
								});
						}),
					root);
			case 'Finished':
				return _elm_lang$core$Native_Utils.update(
					root,
					{finished: true});
			default:
				return A3(
					_concourse$atc$StepTree$updateAt,
					_p43._0,
					_concourse$atc$StepTree$focusRetry(_p43._1),
					root);
		}
	});
var _concourse$atc$StepTree$viewStep = F3(
	function (model, _p44, icon) {
		var _p45 = _p44;
		var _p48 = _p45.state;
		var _p47 = _p45.expanded;
		return A2(
			_elm_lang$html$Html$div,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$classList(
					_elm_lang$core$Native_List.fromArray(
						[
							{ctor: '_Tuple2', _0: 'build-step', _1: true},
							{
							ctor: '_Tuple2',
							_0: 'inactive',
							_1: _elm_lang$core$Basics$not(
								_concourse$atc$StepTree$isActive(_p48))
						},
							{ctor: '_Tuple2', _0: 'first-occurrence', _1: _p45.firstOccurrence}
						]))
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(
					_elm_lang$html$Html$div,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('header'),
							_elm_lang$html$Html_Events$onClick(
							_concourse$atc$StepTree$ToggleStep(_p45.id))
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							A2(_concourse$atc$StepTree$viewStepState, _p48, model.finished),
							_concourse$atc$StepTree$typeIcon(icon),
							_concourse$atc$StepTree$viewVersion(_p45.version),
							A2(
							_elm_lang$html$Html$h3,
							_elm_lang$core$Native_List.fromArray(
								[]),
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html$text(_p45.name)
								])),
							A2(
							_elm_lang$html$Html$div,
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html_Attributes$class('clearfix')
								]),
							_elm_lang$core$Native_List.fromArray(
								[]))
						])),
					A2(
					_elm_lang$html$Html$div,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$classList(
							_elm_lang$core$Native_List.fromArray(
								[
									{ctor: '_Tuple2', _0: 'step-body', _1: true},
									{ctor: '_Tuple2', _0: 'clearfix', _1: true},
									{
									ctor: '_Tuple2',
									_0: 'step-collapsed',
									_1: _elm_lang$core$Basics$not(
										A2(
											_elm_lang$core$Maybe$withDefault,
											_concourse$atc$StepTree$autoExpanded(_p48),
											_p47))
								}
								]))
						]),
					A2(
						_elm_lang$core$Maybe$withDefault,
						_concourse$atc$StepTree$autoExpanded(_p48),
						A2(
							_elm_lang$core$Maybe$map,
							_elm_lang$core$Basics$always(true),
							_p47)) ? _elm_lang$core$Native_List.fromArray(
						[
							_concourse$atc$StepTree$viewMetadata(_p45.metadata),
							_vito$elm_ansi$Ansi_Log$view(_p45.log),
							function () {
							var _p46 = _p45.error;
							if (_p46.ctor === 'Nothing') {
								return A2(
									_elm_lang$html$Html$span,
									_elm_lang$core$Native_List.fromArray(
										[]),
									_elm_lang$core$Native_List.fromArray(
										[]));
							} else {
								return A2(
									_elm_lang$html$Html$span,
									_elm_lang$core$Native_List.fromArray(
										[
											_elm_lang$html$Html_Attributes$class('error')
										]),
									_elm_lang$core$Native_List.fromArray(
										[
											A2(
											_elm_lang$html$Html$pre,
											_elm_lang$core$Native_List.fromArray(
												[]),
											_elm_lang$core$Native_List.fromArray(
												[
													_elm_lang$html$Html$text(_p46._0)
												]))
										]));
							}
						}()
						]) : _elm_lang$core$Native_List.fromArray(
						[]))
				]));
	});
var _concourse$atc$StepTree$viewTree = F2(
	function (model, tree) {
		viewTree:
		while (true) {
			var _p49 = tree;
			switch (_p49.ctor) {
				case 'Task':
					return A3(_concourse$atc$StepTree$viewStep, model, _p49._0, 'fa-terminal');
				case 'Get':
					return A3(_concourse$atc$StepTree$viewStep, model, _p49._0, 'fa-arrow-down');
				case 'DependentGet':
					return A3(_concourse$atc$StepTree$viewStep, model, _p49._0, 'fa-arrow-down');
				case 'Put':
					return A3(_concourse$atc$StepTree$viewStep, model, _p49._0, 'fa-arrow-up');
				case 'Try':
					var _v28 = model,
						_v29 = _p49._0;
					model = _v28;
					tree = _v29;
					continue viewTree;
				case 'Retry':
					var _p53 = _p49._2;
					var _p52 = _p49._1;
					return A2(
						_elm_lang$html$Html$div,
						_elm_lang$core$Native_List.fromArray(
							[
								_elm_lang$html$Html_Attributes$class('retry')
							]),
						_elm_lang$core$Native_List.fromArray(
							[
								A2(
								_elm_lang$html$Html$ul,
								_elm_lang$core$Native_List.fromArray(
									[
										_elm_lang$html$Html_Attributes$class('retry-tabs')
									]),
								_elm_lang$core$Array$toList(
									A2(
										_elm_lang$core$Array$indexedMap,
										A2(_concourse$atc$StepTree$viewTab, _p49._0, _p53),
										_p52))),
								function () {
								var _p50 = A2(_elm_lang$core$Array$get, _p53 - 1, _p52);
								if (_p50.ctor === 'Just') {
									return A2(_concourse$atc$StepTree$viewTree, model, _p50._0);
								} else {
									return _elm_lang$core$Native_Utils.crashCase(
										'StepTree',
										{
											start: {line: 471, column: 11},
											end: {line: 475, column: 60}
										},
										_p50)('impossible (bogus tab selected)');
								}
							}()
							]));
				case 'Timeout':
					var _v31 = model,
						_v32 = _p49._0;
					model = _v31;
					tree = _v32;
					continue viewTree;
				case 'Aggregate':
					return A2(
						_elm_lang$html$Html$div,
						_elm_lang$core$Native_List.fromArray(
							[
								_elm_lang$html$Html_Attributes$class('aggregate')
							]),
						_elm_lang$core$Array$toList(
							A2(
								_elm_lang$core$Array$map,
								_concourse$atc$StepTree$viewSeq(model),
								_p49._0)));
				case 'Do':
					return A2(
						_elm_lang$html$Html$div,
						_elm_lang$core$Native_List.fromArray(
							[
								_elm_lang$html$Html_Attributes$class('do')
							]),
						_elm_lang$core$Array$toList(
							A2(
								_elm_lang$core$Array$map,
								_concourse$atc$StepTree$viewSeq(model),
								_p49._0)));
				case 'OnSuccess':
					return A4(_concourse$atc$StepTree$viewHooked, 'success', model, _p49._0.step, _p49._0.hook);
				case 'OnFailure':
					return A4(_concourse$atc$StepTree$viewHooked, 'failure', model, _p49._0.step, _p49._0.hook);
				default:
					return A4(_concourse$atc$StepTree$viewHooked, 'ensure', model, _p49._0.step, _p49._0.hook);
			}
		}
	});
var _concourse$atc$StepTree$viewHooked = F4(
	function (name, model, step, hook) {
		return A2(
			_elm_lang$html$Html$div,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$class('hooked')
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(
					_elm_lang$html$Html$div,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('step')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							A2(_concourse$atc$StepTree$viewTree, model, step)
						])),
					A2(
					_elm_lang$html$Html$div,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('children')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$html$Html$div,
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html_Attributes$class(
									A2(_elm_lang$core$Basics_ops['++'], 'hook hook-', name))
								]),
							_elm_lang$core$Native_List.fromArray(
								[
									A2(_concourse$atc$StepTree$viewTree, model, hook)
								]))
						]))
				]));
	});
var _concourse$atc$StepTree$viewSeq = F2(
	function (model, tree) {
		return A2(
			_elm_lang$html$Html$div,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$class('seq')
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(_concourse$atc$StepTree$viewTree, model, tree)
				]));
	});
var _concourse$atc$StepTree$view = function (model) {
	return A2(_concourse$atc$StepTree$viewTree, model, model.tree);
};

var _concourse$atc$BuildOutput$viewLoginButton = function (build) {
	return A2(
		_elm_lang$html$Html$form,
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html_Attributes$class('build-login'),
				_elm_lang$html$Html_Attributes$method('get'),
				_elm_lang$html$Html_Attributes$action('/login')
			]),
		_elm_lang$core$Native_List.fromArray(
			[
				A2(
				_elm_lang$html$Html$input,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$type$('submit'),
						_elm_lang$html$Html_Attributes$value('log in to view')
					]),
				_elm_lang$core$Native_List.fromArray(
					[])),
				A2(
				_elm_lang$html$Html$input,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$type$('hidden'),
						_elm_lang$html$Html_Attributes$name('redirect'),
						_elm_lang$html$Html_Attributes$value(
						_concourse$atc$Concourse_Build$url(build))
					]),
				_elm_lang$core$Native_List.fromArray(
					[]))
			]));
};
var _concourse$atc$BuildOutput$viewErrors = function (errors) {
	var _p0 = errors;
	if (_p0.ctor === 'Nothing') {
		return A2(
			_elm_lang$html$Html$div,
			_elm_lang$core$Native_List.fromArray(
				[]),
			_elm_lang$core$Native_List.fromArray(
				[]));
	} else {
		return A2(
			_elm_lang$html$Html$div,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$class('build-step')
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(
					_elm_lang$html$Html$div,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('header')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$html$Html$i,
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html_Attributes$class('left fa fa-fw fa-exclamation-triangle')
								]),
							_elm_lang$core$Native_List.fromArray(
								[])),
							A2(
							_elm_lang$html$Html$h3,
							_elm_lang$core$Native_List.fromArray(
								[]),
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html$text('error')
								]))
						])),
					A2(
					_elm_lang$html$Html$div,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('step-body build-errors-body')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							_vito$elm_ansi$Ansi_Log$view(_p0._0)
						]))
				]));
	}
};
var _concourse$atc$BuildOutput$setStepState = F2(
	function (state, tree) {
		return A2(
			_concourse$atc$StepTree$map,
			function (step) {
				return _elm_lang$core$Native_Utils.update(
					step,
					{state: state});
			},
			tree);
	});
var _concourse$atc$BuildOutput$setResourceInfo = F3(
	function (version, metadata, tree) {
		return A2(
			_concourse$atc$StepTree$map,
			function (step) {
				return _elm_lang$core$Native_Utils.update(
					step,
					{
						version: _elm_lang$core$Maybe$Just(version),
						metadata: metadata
					});
			},
			tree);
	});
var _concourse$atc$BuildOutput$finishStep = F2(
	function (exitStatus, tree) {
		var stepState = _elm_lang$core$Native_Utils.eq(exitStatus, 0) ? _concourse$atc$StepTree$StepStateSucceeded : _concourse$atc$StepTree$StepStateFailed;
		return A2(_concourse$atc$BuildOutput$setStepState, stepState, tree);
	});
var _concourse$atc$BuildOutput$setStepError = F2(
	function (message, tree) {
		return A2(
			_concourse$atc$StepTree$map,
			function (step) {
				return _elm_lang$core$Native_Utils.update(
					step,
					{
						state: _concourse$atc$StepTree$StepStateErrored,
						error: _elm_lang$core$Maybe$Just(message)
					});
			},
			tree);
	});
var _concourse$atc$BuildOutput$appendStepLog = F2(
	function (output, tree) {
		return A2(
			_concourse$atc$StepTree$map,
			function (step) {
				return _elm_lang$core$Native_Utils.update(
					step,
					{
						log: A2(_vito$elm_ansi$Ansi_Log$update, output, step.log)
					});
			},
			tree);
	});
var _concourse$atc$BuildOutput$setRunning = _concourse$atc$BuildOutput$setStepState(_concourse$atc$StepTree$StepStateRunning);
var _concourse$atc$BuildOutput$updateStep = F3(
	function (id, update, model) {
		return _elm_lang$core$Native_Utils.update(
			model,
			{
				steps: A2(
					_elm_lang$core$Maybe$map,
					A2(_concourse$atc$StepTree$updateAt, id, update),
					model.steps)
			});
	});
var _concourse$atc$BuildOutput$handleEvent = F2(
	function (event, model) {
		var _p1 = event;
		switch (_p1.ctor) {
			case 'Log':
				return {
					ctor: '_Tuple2',
					_0: A3(
						_concourse$atc$BuildOutput$updateStep,
						_p1._0.id,
						function (_p2) {
							return _concourse$atc$BuildOutput$setRunning(
								A2(_concourse$atc$BuildOutput$appendStepLog, _p1._1, _p2));
						},
						model),
					_1: _elm_lang$core$Platform_Cmd$none
				};
			case 'Error':
				return {
					ctor: '_Tuple2',
					_0: A3(
						_concourse$atc$BuildOutput$updateStep,
						_p1._0.id,
						_concourse$atc$BuildOutput$setStepError(_p1._1),
						model),
					_1: _elm_lang$core$Platform_Cmd$none
				};
			case 'InitializeTask':
				return {
					ctor: '_Tuple2',
					_0: A3(_concourse$atc$BuildOutput$updateStep, _p1._0.id, _concourse$atc$BuildOutput$setRunning, model),
					_1: _elm_lang$core$Platform_Cmd$none
				};
			case 'StartTask':
				return {
					ctor: '_Tuple2',
					_0: A3(_concourse$atc$BuildOutput$updateStep, _p1._0.id, _concourse$atc$BuildOutput$setRunning, model),
					_1: _elm_lang$core$Platform_Cmd$none
				};
			case 'FinishTask':
				return {
					ctor: '_Tuple2',
					_0: A3(
						_concourse$atc$BuildOutput$updateStep,
						_p1._0.id,
						_concourse$atc$BuildOutput$finishStep(_p1._1),
						model),
					_1: _elm_lang$core$Platform_Cmd$none
				};
			case 'InitializeGet':
				return {
					ctor: '_Tuple2',
					_0: A3(_concourse$atc$BuildOutput$updateStep, _p1._0.id, _concourse$atc$BuildOutput$setRunning, model),
					_1: _elm_lang$core$Platform_Cmd$none
				};
			case 'FinishGet':
				return {
					ctor: '_Tuple2',
					_0: A3(
						_concourse$atc$BuildOutput$updateStep,
						_p1._0.id,
						function (_p3) {
							return A2(
								_concourse$atc$BuildOutput$finishStep,
								_p1._1,
								A3(_concourse$atc$BuildOutput$setResourceInfo, _p1._2, _p1._3, _p3));
						},
						model),
					_1: _elm_lang$core$Platform_Cmd$none
				};
			case 'InitializePut':
				return {
					ctor: '_Tuple2',
					_0: A3(_concourse$atc$BuildOutput$updateStep, _p1._0.id, _concourse$atc$BuildOutput$setRunning, model),
					_1: _elm_lang$core$Platform_Cmd$none
				};
			case 'FinishPut':
				return {
					ctor: '_Tuple2',
					_0: A3(
						_concourse$atc$BuildOutput$updateStep,
						_p1._0.id,
						function (_p4) {
							return A2(
								_concourse$atc$BuildOutput$finishStep,
								_p1._1,
								A3(_concourse$atc$BuildOutput$setResourceInfo, _p1._2, _p1._3, _p4));
						},
						model),
					_1: _elm_lang$core$Platform_Cmd$none
				};
			case 'BuildStatus':
				var _p6 = _p1._0;
				return {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Native_Utils.update(
						model,
						{
							steps: _elm_lang$core$Basics$not(
								_concourse$atc$Concourse_BuildStatus$isRunning(_p6)) ? A2(
								_elm_lang$core$Maybe$map,
								_concourse$atc$StepTree$update(_concourse$atc$StepTree$Finished),
								model.steps) : model.steps,
							build: function () {
								var _p5 = {ctor: '_Tuple3', _0: _p6, _1: model.build, _2: model.build.duration};
								switch (_p5._0.ctor) {
									case 'Started':
										return model.build;
									case 'Pending':
										return model.build;
									default:
										return _elm_lang$core$Native_Utils.update(
											_p5._1,
											{
												status: _p6,
												duration: _elm_lang$core$Native_Utils.update(
													_p5._2,
													{
														finishedAt: _elm_lang$core$Maybe$Just(_p1._1)
													})
											});
								}
							}()
						}),
					_1: _elm_lang$core$Platform_Cmd$none
				};
			default:
				return {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Native_Utils.update(
						model,
						{
							errors: _elm_lang$core$Maybe$Just(
								A2(
									_vito$elm_ansi$Ansi_Log$update,
									_p1._0,
									A2(
										_elm_lang$core$Maybe$withDefault,
										_vito$elm_ansi$Ansi_Log$init(_vito$elm_ansi$Ansi_Log$Cooked),
										model.errors)))
						}),
					_1: _elm_lang$core$Platform_Cmd$none
				};
		}
	});
var _concourse$atc$BuildOutput$Model = F6(
	function (a, b, c, d, e, f) {
		return {build: a, steps: b, errors: c, state: d, eventSourceOpened: e, events: f};
	});
var _concourse$atc$BuildOutput$LoginRequired = {ctor: 'LoginRequired'};
var _concourse$atc$BuildOutput$StepsComplete = {ctor: 'StepsComplete'};
var _concourse$atc$BuildOutput$handleEventsAction = F2(
	function (action, model) {
		var _p7 = action;
		switch (_p7.ctor) {
			case 'Opened':
				return {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Native_Utils.update(
						model,
						{eventSourceOpened: true}),
					_1: _elm_lang$core$Platform_Cmd$none
				};
			case 'Errored':
				return model.eventSourceOpened ? {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none} : {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Native_Utils.update(
						model,
						{state: _concourse$atc$BuildOutput$LoginRequired}),
					_1: _elm_lang$core$Platform_Cmd$none
				};
			case 'Event':
				if (_p7._0.ctor === 'Ok') {
					return A2(_concourse$atc$BuildOutput$handleEvent, _p7._0._0, model);
				} else {
					return {
						ctor: '_Tuple2',
						_0: model,
						_1: A2(_elm_lang$core$Debug$log, _p7._0._0, _elm_lang$core$Platform_Cmd$none)
					};
				}
			default:
				return {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Native_Utils.update(
						model,
						{state: _concourse$atc$BuildOutput$StepsComplete, events: _elm_lang$core$Platform_Sub$none}),
					_1: _elm_lang$core$Platform_Cmd$none
				};
		}
	});
var _concourse$atc$BuildOutput$StepsLiveUpdating = {ctor: 'StepsLiveUpdating'};
var _concourse$atc$BuildOutput$StepsLoading = {ctor: 'StepsLoading'};
var _concourse$atc$BuildOutput$StepTreeAction = function (a) {
	return {ctor: 'StepTreeAction', _0: a};
};
var _concourse$atc$BuildOutput$viewStepTree = F3(
	function (build, steps, state) {
		var _p8 = {ctor: '_Tuple2', _0: state, _1: steps};
		_v4_4:
		do {
			switch (_p8._0.ctor) {
				case 'StepsLoading':
					return _concourse$atc$LoadingIndicator$view;
				case 'LoginRequired':
					return _concourse$atc$BuildOutput$viewLoginButton(build);
				case 'StepsLiveUpdating':
					if (_p8._1.ctor === 'Just') {
						return A2(
							_elm_lang$html$Html_App$map,
							_concourse$atc$BuildOutput$StepTreeAction,
							_concourse$atc$StepTree$view(_p8._1._0));
					} else {
						break _v4_4;
					}
				default:
					if (_p8._1.ctor === 'Just') {
						return A2(
							_elm_lang$html$Html_App$map,
							_concourse$atc$BuildOutput$StepTreeAction,
							_concourse$atc$StepTree$view(_p8._1._0));
					} else {
						break _v4_4;
					}
			}
		} while(false);
		return A2(
			_elm_lang$html$Html$div,
			_elm_lang$core$Native_List.fromArray(
				[]),
			_elm_lang$core$Native_List.fromArray(
				[]));
	});
var _concourse$atc$BuildOutput$view = function (_p9) {
	var _p10 = _p9;
	return A2(
		_elm_lang$html$Html$div,
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html_Attributes$class('steps')
			]),
		_elm_lang$core$Native_List.fromArray(
			[
				_concourse$atc$BuildOutput$viewErrors(_p10.errors),
				A3(_concourse$atc$BuildOutput$viewStepTree, _p10.build, _p10.steps, _p10.state)
			]));
};
var _concourse$atc$BuildOutput$BuildEventsAction = function (a) {
	return {ctor: 'BuildEventsAction', _0: a};
};
var _concourse$atc$BuildOutput$subscribeToEvents = function (build) {
	return A2(
		_elm_lang$core$Platform_Sub$map,
		_concourse$atc$BuildOutput$BuildEventsAction,
		_concourse$atc$Concourse_BuildEvents$subscribe(build));
};
var _concourse$atc$BuildOutput$update = F2(
	function (action, model) {
		var _p11 = action;
		switch (_p11.ctor) {
			case 'Noop':
				return {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none};
			case 'PlanAndResourcesFetched':
				if (_p11._0.ctor === 'Err') {
					if ((_p11._0._0.ctor === 'BadResponse') && (_p11._0._0._0 === 404)) {
						return {
							ctor: '_Tuple2',
							_0: _elm_lang$core$Native_Utils.update(
								model,
								{
									events: _concourse$atc$BuildOutput$subscribeToEvents(model.build.id)
								}),
							_1: _elm_lang$core$Platform_Cmd$none
						};
					} else {
						return A2(
							_elm_lang$core$Debug$log,
							A2(
								_elm_lang$core$Basics_ops['++'],
								'failed to fetch plan: ',
								_elm_lang$core$Basics$toString(_p11._0._0)),
							{ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none});
					}
				} else {
					return {
						ctor: '_Tuple2',
						_0: _elm_lang$core$Native_Utils.update(
							model,
							{
								steps: _elm_lang$core$Maybe$Just(
									A2(_concourse$atc$StepTree$init, _p11._0._0._1, _p11._0._0._0)),
								events: _concourse$atc$BuildOutput$subscribeToEvents(model.build.id)
							}),
						_1: _elm_lang$core$Platform_Cmd$none
					};
				}
			case 'BuildEventsAction':
				return A2(_concourse$atc$BuildOutput$handleEventsAction, _p11._0, model);
			default:
				return {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Native_Utils.update(
						model,
						{
							steps: A2(
								_elm_lang$core$Maybe$map,
								_concourse$atc$StepTree$update(_p11._0),
								model.steps)
						}),
					_1: _elm_lang$core$Platform_Cmd$none
				};
		}
	});
var _concourse$atc$BuildOutput$PlanAndResourcesFetched = function (a) {
	return {ctor: 'PlanAndResourcesFetched', _0: a};
};
var _concourse$atc$BuildOutput$fetchBuildPlanAndResources = function (buildId) {
	return function (_p12) {
		return A2(
			_elm_lang$core$Platform_Cmd$map,
			_concourse$atc$BuildOutput$PlanAndResourcesFetched,
			A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _p12));
	}(
		A3(
			_elm_lang$core$Task$map2,
			F2(
				function (v0, v1) {
					return {ctor: '_Tuple2', _0: v0, _1: v1};
				}),
			_concourse$atc$Concourse_BuildPlan$fetch(buildId),
			_concourse$atc$Concourse_BuildResources$fetch(buildId)));
};
var _concourse$atc$BuildOutput$fetchBuildPlan = function (buildId) {
	return function (_p13) {
		return A2(
			_elm_lang$core$Platform_Cmd$map,
			_concourse$atc$BuildOutput$PlanAndResourcesFetched,
			A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _p13));
	}(
		A2(
			_elm_lang$core$Task$map,
			A2(
				_elm_lang$core$Basics$flip,
				F2(
					function (v0, v1) {
						return {ctor: '_Tuple2', _0: v0, _1: v1};
					}),
				_concourse$atc$Concourse_BuildResources$empty),
			_concourse$atc$Concourse_BuildPlan$fetch(buildId)));
};
var _concourse$atc$BuildOutput$init = function (build) {
	var outputState = _concourse$atc$Concourse_BuildStatus$isRunning(build.status) ? _concourse$atc$BuildOutput$StepsLiveUpdating : _concourse$atc$BuildOutput$StepsLoading;
	var model = {build: build, steps: _elm_lang$core$Maybe$Nothing, errors: _elm_lang$core$Maybe$Nothing, state: outputState, events: _elm_lang$core$Platform_Sub$none, eventSourceOpened: false};
	var fetch = (!_elm_lang$core$Native_Utils.eq(build.job, _elm_lang$core$Maybe$Nothing)) ? _concourse$atc$BuildOutput$fetchBuildPlanAndResources(model.build.id) : _concourse$atc$BuildOutput$fetchBuildPlan(model.build.id);
	return {ctor: '_Tuple2', _0: model, _1: fetch};
};
var _concourse$atc$BuildOutput$Noop = {ctor: 'Noop'};

var _concourse$atc$Concourse_BuildPrep$replaceMaybeWithEmptyDict = function (maybeDict) {
	var _p0 = maybeDict;
	if (_p0.ctor === 'Just') {
		return _p0._0;
	} else {
		return _elm_lang$core$Dict$empty;
	}
};
var _concourse$atc$Concourse_BuildPrep$BuildPrep = F6(
	function (a, b, c, d, e, f) {
		return {pausedPipeline: a, pausedJob: b, maxRunningBuilds: c, inputs: d, inputsSatisfied: e, missingInputReasons: f};
	});
var _concourse$atc$Concourse_BuildPrep$NotBlocking = {ctor: 'NotBlocking'};
var _concourse$atc$Concourse_BuildPrep$Blocking = {ctor: 'Blocking'};
var _concourse$atc$Concourse_BuildPrep$Unknown = {ctor: 'Unknown'};
var _concourse$atc$Concourse_BuildPrep$decodeStatus = A2(
	_elm_lang$core$Json_Decode$customDecoder,
	_elm_lang$core$Json_Decode$string,
	function (status) {
		var _p1 = status;
		switch (_p1) {
			case 'unknown':
				return _elm_lang$core$Result$Ok(_concourse$atc$Concourse_BuildPrep$Unknown);
			case 'blocking':
				return _elm_lang$core$Result$Ok(_concourse$atc$Concourse_BuildPrep$Blocking);
			case 'not_blocking':
				return _elm_lang$core$Result$Ok(_concourse$atc$Concourse_BuildPrep$NotBlocking);
			default:
				return _elm_lang$core$Result$Err(
					A2(_elm_lang$core$Basics_ops['++'], 'unknown build preparation status: ', _p1));
		}
	});
var _concourse$atc$Concourse_BuildPrep$decode = A7(
	_elm_lang$core$Json_Decode$object6,
	_concourse$atc$Concourse_BuildPrep$BuildPrep,
	A2(_elm_lang$core$Json_Decode_ops[':='], 'paused_pipeline', _concourse$atc$Concourse_BuildPrep$decodeStatus),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'paused_job', _concourse$atc$Concourse_BuildPrep$decodeStatus),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'max_running_builds', _concourse$atc$Concourse_BuildPrep$decodeStatus),
	A2(
		_elm_lang$core$Json_Decode_ops[':='],
		'inputs',
		_elm_lang$core$Json_Decode$dict(_concourse$atc$Concourse_BuildPrep$decodeStatus)),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'inputs_satisfied', _concourse$atc$Concourse_BuildPrep$decodeStatus),
	A2(
		_elm_lang$core$Json_Decode$map,
		_concourse$atc$Concourse_BuildPrep$replaceMaybeWithEmptyDict,
		_elm_lang$core$Json_Decode$maybe(
			A2(
				_elm_lang$core$Json_Decode_ops[':='],
				'missing_input_reasons',
				_elm_lang$core$Json_Decode$dict(_elm_lang$core$Json_Decode$string)))));
var _concourse$atc$Concourse_BuildPrep$fetch = function (buildId) {
	return A2(
		_evancz$elm_http$Http$get,
		_concourse$atc$Concourse_BuildPrep$decode,
		A2(
			_elm_lang$core$Basics_ops['++'],
			'/api/v1/builds/',
			A2(
				_elm_lang$core$Basics_ops['++'],
				_elm_lang$core$Basics$toString(buildId),
				'/preparation')));
};

var _concourse$atc$Duration$format = function (duration) {
	var seconds = _elm_lang$core$Basics$truncate(duration / 1000);
	var remainingSeconds = A2(_elm_lang$core$Basics$rem, seconds, 60);
	var minutes = (seconds / 60) | 0;
	var remainingMinutes = A2(_elm_lang$core$Basics$rem, minutes, 60);
	var hours = (minutes / 60) | 0;
	var remainingHours = A2(_elm_lang$core$Basics$rem, hours, 24);
	var days = (hours / 24) | 0;
	var _p0 = {ctor: '_Tuple4', _0: days, _1: remainingHours, _2: remainingMinutes, _3: remainingSeconds};
	if (_p0._0 === 0) {
		if (_p0._1 === 0) {
			if (_p0._2 === 0) {
				return A2(
					_elm_lang$core$Basics_ops['++'],
					_elm_lang$core$Basics$toString(_p0._3),
					's');
			} else {
				return A2(
					_elm_lang$core$Basics_ops['++'],
					_elm_lang$core$Basics$toString(_p0._2),
					A2(
						_elm_lang$core$Basics_ops['++'],
						'm ',
						A2(
							_elm_lang$core$Basics_ops['++'],
							_elm_lang$core$Basics$toString(_p0._3),
							's')));
			}
		} else {
			return A2(
				_elm_lang$core$Basics_ops['++'],
				_elm_lang$core$Basics$toString(_p0._1),
				A2(
					_elm_lang$core$Basics_ops['++'],
					'h ',
					A2(
						_elm_lang$core$Basics_ops['++'],
						_elm_lang$core$Basics$toString(_p0._2),
						'm')));
		}
	} else {
		return A2(
			_elm_lang$core$Basics_ops['++'],
			_elm_lang$core$Basics$toString(_p0._0),
			A2(
				_elm_lang$core$Basics_ops['++'],
				'd ',
				A2(
					_elm_lang$core$Basics_ops['++'],
					_elm_lang$core$Basics$toString(_p0._1),
					'h')));
	}
};
var _concourse$atc$Duration$between = F2(
	function (a, b) {
		return b - a;
	});

var _concourse$atc$BuildDuration$pendingLabel = function (label) {
	return A2(
		_elm_lang$html$Html$tr,
		_elm_lang$core$Native_List.fromArray(
			[]),
		_elm_lang$core$Native_List.fromArray(
			[
				A2(
				_elm_lang$html$Html$td,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('dict-key')
					]),
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html$text(label)
					])),
				A2(
				_elm_lang$html$Html$td,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('dict-value')
					]),
				_elm_lang$core$Native_List.fromArray(
					[]))
			]));
};
var _concourse$atc$BuildDuration$labeledDuration = F2(
	function (label, duration) {
		return A2(
			_elm_lang$html$Html$tr,
			_elm_lang$core$Native_List.fromArray(
				[]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(
					_elm_lang$html$Html$td,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('dict-key')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html$text(label)
						])),
					A2(
					_elm_lang$html$Html$td,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('dict-value')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$html$Html$span,
							_elm_lang$core$Native_List.fromArray(
								[]),
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html$text(
									_concourse$atc$Duration$format(duration))
								]))
						]))
				]));
	});
var _concourse$atc$BuildDuration$labeledRelativeDate = F3(
	function (label, now, date) {
		var ago = A2(
			_concourse$atc$Duration$between,
			_elm_lang$core$Date$toTime(date),
			now);
		return A2(
			_elm_lang$html$Html$tr,
			_elm_lang$core$Native_List.fromArray(
				[]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(
					_elm_lang$html$Html$td,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('dict-key')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html$text(label)
						])),
					A2(
					_elm_lang$html$Html$td,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$title(
							A2(_mgold$elm_date_format$Date_Format$format, '%b %d %Y %I:%M:%S %p', date)),
							_elm_lang$html$Html_Attributes$class('dict-value')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$html$Html$span,
							_elm_lang$core$Native_List.fromArray(
								[]),
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html$text(
									A2(
										_elm_lang$core$Basics_ops['++'],
										_concourse$atc$Duration$format(ago),
										' ago'))
								]))
						]))
				]));
	});
var _concourse$atc$BuildDuration$view = F2(
	function (duration, now) {
		return A2(
			_elm_lang$html$Html$table,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$class('dictionary build-duration')
				]),
			function () {
				var _p0 = {ctor: '_Tuple2', _0: duration.startedAt, _1: duration.finishedAt};
				if (_p0._0.ctor === 'Nothing') {
					if (_p0._1.ctor === 'Nothing') {
						return _elm_lang$core$Native_List.fromArray(
							[
								_concourse$atc$BuildDuration$pendingLabel('pending')
							]);
					} else {
						return _elm_lang$core$Native_List.fromArray(
							[
								A3(_concourse$atc$BuildDuration$labeledRelativeDate, 'finished', now, _p0._1._0)
							]);
					}
				} else {
					if (_p0._1.ctor === 'Nothing') {
						return _elm_lang$core$Native_List.fromArray(
							[
								A3(_concourse$atc$BuildDuration$labeledRelativeDate, 'started', now, _p0._0._0)
							]);
					} else {
						var _p2 = _p0._0._0;
						var _p1 = _p0._1._0;
						var durationElmIssue = A2(
							_concourse$atc$Duration$between,
							_elm_lang$core$Date$toTime(_p2),
							_elm_lang$core$Date$toTime(_p1));
						return _elm_lang$core$Native_List.fromArray(
							[
								A3(_concourse$atc$BuildDuration$labeledRelativeDate, 'started', now, _p2),
								A3(_concourse$atc$BuildDuration$labeledRelativeDate, 'finished', now, _p1),
								A2(_concourse$atc$BuildDuration$labeledDuration, 'duration', durationElmIssue)
							]);
					}
				}
			}());
	});

var _concourse$atc$Native_Redirect = function() {
  function to(url) {
    return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback) {
      window.location = url;
      callback(_elm_lang$core$Native_Scheduler.succeed(_elm_lang$core$Native_Utils.Tuple0));
    });
  }

  return {
    to: to
  };
}();

var _concourse$atc$Redirect$to = _concourse$atc$Native_Redirect.to;

var _elm_lang$core$Set$foldr = F3(
	function (f, b, _p0) {
		var _p1 = _p0;
		return A3(
			_elm_lang$core$Dict$foldr,
			F3(
				function (k, _p2, b) {
					return A2(f, k, b);
				}),
			b,
			_p1._0);
	});
var _elm_lang$core$Set$foldl = F3(
	function (f, b, _p3) {
		var _p4 = _p3;
		return A3(
			_elm_lang$core$Dict$foldl,
			F3(
				function (k, _p5, b) {
					return A2(f, k, b);
				}),
			b,
			_p4._0);
	});
var _elm_lang$core$Set$toList = function (_p6) {
	var _p7 = _p6;
	return _elm_lang$core$Dict$keys(_p7._0);
};
var _elm_lang$core$Set$size = function (_p8) {
	var _p9 = _p8;
	return _elm_lang$core$Dict$size(_p9._0);
};
var _elm_lang$core$Set$member = F2(
	function (k, _p10) {
		var _p11 = _p10;
		return A2(_elm_lang$core$Dict$member, k, _p11._0);
	});
var _elm_lang$core$Set$isEmpty = function (_p12) {
	var _p13 = _p12;
	return _elm_lang$core$Dict$isEmpty(_p13._0);
};
var _elm_lang$core$Set$Set_elm_builtin = function (a) {
	return {ctor: 'Set_elm_builtin', _0: a};
};
var _elm_lang$core$Set$empty = _elm_lang$core$Set$Set_elm_builtin(_elm_lang$core$Dict$empty);
var _elm_lang$core$Set$singleton = function (k) {
	return _elm_lang$core$Set$Set_elm_builtin(
		A2(
			_elm_lang$core$Dict$singleton,
			k,
			{ctor: '_Tuple0'}));
};
var _elm_lang$core$Set$insert = F2(
	function (k, _p14) {
		var _p15 = _p14;
		return _elm_lang$core$Set$Set_elm_builtin(
			A3(
				_elm_lang$core$Dict$insert,
				k,
				{ctor: '_Tuple0'},
				_p15._0));
	});
var _elm_lang$core$Set$fromList = function (xs) {
	return A3(_elm_lang$core$List$foldl, _elm_lang$core$Set$insert, _elm_lang$core$Set$empty, xs);
};
var _elm_lang$core$Set$map = F2(
	function (f, s) {
		return _elm_lang$core$Set$fromList(
			A2(
				_elm_lang$core$List$map,
				f,
				_elm_lang$core$Set$toList(s)));
	});
var _elm_lang$core$Set$remove = F2(
	function (k, _p16) {
		var _p17 = _p16;
		return _elm_lang$core$Set$Set_elm_builtin(
			A2(_elm_lang$core$Dict$remove, k, _p17._0));
	});
var _elm_lang$core$Set$union = F2(
	function (_p19, _p18) {
		var _p20 = _p19;
		var _p21 = _p18;
		return _elm_lang$core$Set$Set_elm_builtin(
			A2(_elm_lang$core$Dict$union, _p20._0, _p21._0));
	});
var _elm_lang$core$Set$intersect = F2(
	function (_p23, _p22) {
		var _p24 = _p23;
		var _p25 = _p22;
		return _elm_lang$core$Set$Set_elm_builtin(
			A2(_elm_lang$core$Dict$intersect, _p24._0, _p25._0));
	});
var _elm_lang$core$Set$diff = F2(
	function (_p27, _p26) {
		var _p28 = _p27;
		var _p29 = _p26;
		return _elm_lang$core$Set$Set_elm_builtin(
			A2(_elm_lang$core$Dict$diff, _p28._0, _p29._0));
	});
var _elm_lang$core$Set$filter = F2(
	function (p, _p30) {
		var _p31 = _p30;
		return _elm_lang$core$Set$Set_elm_builtin(
			A2(
				_elm_lang$core$Dict$filter,
				F2(
					function (k, _p32) {
						return p(k);
					}),
				_p31._0));
	});
var _elm_lang$core$Set$partition = F2(
	function (p, _p33) {
		var _p34 = _p33;
		var _p35 = A2(
			_elm_lang$core$Dict$partition,
			F2(
				function (k, _p36) {
					return p(k);
				}),
			_p34._0);
		var p1 = _p35._0;
		var p2 = _p35._1;
		return {
			ctor: '_Tuple2',
			_0: _elm_lang$core$Set$Set_elm_builtin(p1),
			_1: _elm_lang$core$Set$Set_elm_builtin(p2)
		};
	});

var _elm_community$json_extra$Json_Decode_Extra$lazy = function (getDecoder) {
	return A2(
		_elm_lang$core$Json_Decode$customDecoder,
		_elm_lang$core$Json_Decode$value,
		function (rawValue) {
			return A2(
				_elm_lang$core$Json_Decode$decodeValue,
				getDecoder(
					{ctor: '_Tuple0'}),
				rawValue);
		});
};
var _elm_community$json_extra$Json_Decode_Extra$maybeNull = function (decoder) {
	return _elm_lang$core$Json_Decode$oneOf(
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$core$Json_Decode$null(_elm_lang$core$Maybe$Nothing),
				A2(_elm_lang$core$Json_Decode$map, _elm_lang$core$Maybe$Just, decoder)
			]));
};
var _elm_community$json_extra$Json_Decode_Extra$withDefault = F2(
	function (fallback, decoder) {
		return A2(
			_elm_lang$core$Json_Decode$andThen,
			_elm_lang$core$Json_Decode$maybe(decoder),
			function (_p0) {
				return _elm_lang$core$Json_Decode$succeed(
					A2(_elm_lang$core$Maybe$withDefault, fallback, _p0));
			});
	});
var _elm_community$json_extra$Json_Decode_Extra$decodeDictFromTuples = F2(
	function (keyDecoder, tuples) {
		var _p1 = tuples;
		if (_p1.ctor === '[]') {
			return _elm_lang$core$Json_Decode$succeed(_elm_lang$core$Dict$empty);
		} else {
			var _p2 = A2(_elm_lang$core$Json_Decode$decodeString, keyDecoder, _p1._0._0);
			if (_p2.ctor === 'Ok') {
				return A2(
					_elm_lang$core$Json_Decode$andThen,
					A2(_elm_community$json_extra$Json_Decode_Extra$decodeDictFromTuples, keyDecoder, _p1._1),
					function (_p3) {
						return _elm_lang$core$Json_Decode$succeed(
							A3(_elm_lang$core$Dict$insert, _p2._0, _p1._0._1, _p3));
					});
			} else {
				return _elm_lang$core$Json_Decode$fail(_p2._0);
			}
		}
	});
var _elm_community$json_extra$Json_Decode_Extra$dict2 = F2(
	function (keyDecoder, valueDecoder) {
		return A2(
			_elm_lang$core$Json_Decode$andThen,
			_elm_lang$core$Json_Decode$dict(valueDecoder),
			function (_p4) {
				return A2(
					_elm_community$json_extra$Json_Decode_Extra$decodeDictFromTuples,
					keyDecoder,
					_elm_lang$core$Dict$toList(_p4));
			});
	});
var _elm_community$json_extra$Json_Decode_Extra$set = function (decoder) {
	return A2(
		_elm_lang$core$Json_Decode$andThen,
		_elm_lang$core$Json_Decode$list(decoder),
		function (_p5) {
			return _elm_lang$core$Json_Decode$succeed(
				_elm_lang$core$Set$fromList(_p5));
		});
};
var _elm_community$json_extra$Json_Decode_Extra$date = A2(_elm_lang$core$Json_Decode$customDecoder, _elm_lang$core$Json_Decode$string, _elm_lang$core$Date$fromString);
var _elm_community$json_extra$Json_Decode_Extra$apply = _elm_lang$core$Json_Decode$object2(
	F2(
		function (x, y) {
			return x(y);
		}));
var _elm_community$json_extra$Json_Decode_Extra_ops = _elm_community$json_extra$Json_Decode_Extra_ops || {};
_elm_community$json_extra$Json_Decode_Extra_ops['|:'] = _elm_community$json_extra$Json_Decode_Extra$apply;

var _concourse$atc$Concourse_Job$promoteHttpError = function (rawError) {
	var _p0 = rawError;
	if (_p0.ctor === 'RawTimeout') {
		return _evancz$elm_http$Http$Timeout;
	} else {
		return _evancz$elm_http$Http$NetworkError;
	}
};
var _concourse$atc$Concourse_Job$handleResponse = function (response) {
	return ((_elm_lang$core$Native_Utils.cmp(200, response.status) < 1) && (_elm_lang$core$Native_Utils.cmp(response.status, 300) < 0)) ? _elm_lang$core$Task$succeed(
		{ctor: '_Tuple0'}) : _elm_lang$core$Task$fail(
		A2(_evancz$elm_http$Http$BadResponse, response.status, response.statusText));
};
var _concourse$atc$Concourse_Job$pauseUnpause = F2(
	function (pause, jobInfo) {
		var action = pause ? 'pause' : 'unpause';
		var put = A2(
			_evancz$elm_http$Http$send,
			_evancz$elm_http$Http$defaultSettings,
			{
				verb: 'PUT',
				headers: _elm_lang$core$Native_List.fromArray(
					[]),
				url: A2(
					_elm_lang$core$Basics_ops['++'],
					'/api/v1/teams/',
					A2(
						_elm_lang$core$Basics_ops['++'],
						jobInfo.teamName,
						A2(
							_elm_lang$core$Basics_ops['++'],
							'/pipelines/',
							A2(
								_elm_lang$core$Basics_ops['++'],
								jobInfo.pipelineName,
								A2(
									_elm_lang$core$Basics_ops['++'],
									'/jobs/',
									A2(
										_elm_lang$core$Basics_ops['++'],
										jobInfo.name,
										A2(_elm_lang$core$Basics_ops['++'], '/', action))))))),
				body: _evancz$elm_http$Http$empty
			});
		return A2(
			_elm_lang$core$Task$andThen,
			A2(_elm_lang$core$Task$mapError, _concourse$atc$Concourse_Job$promoteHttpError, put),
			_concourse$atc$Concourse_Job$handleResponse);
	});
var _concourse$atc$Concourse_Job$unpause = function (jobInfo) {
	return A2(_concourse$atc$Concourse_Job$pauseUnpause, false, jobInfo);
};
var _concourse$atc$Concourse_Job$pause = function (jobInfo) {
	return A2(_concourse$atc$Concourse_Job$pauseUnpause, true, jobInfo);
};
var _concourse$atc$Concourse_Job$optional = function ($default) {
	return _elm_lang$core$Json_Decode$map(
		_elm_lang$core$Maybe$withDefault($default));
};
var _concourse$atc$Concourse_Job$Job = function (a) {
	return function (b) {
		return function (c) {
			return function (d) {
				return function (e) {
					return function (f) {
						return function (g) {
							return function (h) {
								return function (i) {
									return function (j) {
										return function (k) {
											return {teamName: a, pipelineName: b, name: c, url: d, nextBuild: e, finishedBuild: f, paused: g, disableManualTrigger: h, inputs: i, outputs: j, groups: k};
										};
									};
								};
							};
						};
					};
				};
			};
		};
	};
};
var _concourse$atc$Concourse_Job$Input = F4(
	function (a, b, c, d) {
		return {name: a, resource: b, passed: c, trigger: d};
	});
var _concourse$atc$Concourse_Job$decodeInput = A5(
	_elm_lang$core$Json_Decode$object4,
	_concourse$atc$Concourse_Job$Input,
	A2(_elm_lang$core$Json_Decode_ops[':='], 'name', _elm_lang$core$Json_Decode$string),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'resource', _elm_lang$core$Json_Decode$string),
	A2(
		_concourse$atc$Concourse_Job$optional,
		_elm_lang$core$Native_List.fromArray(
			[]),
		_elm_lang$core$Json_Decode$maybe(
			A2(
				_elm_lang$core$Json_Decode_ops[':='],
				'passed',
				_elm_lang$core$Json_Decode$list(_elm_lang$core$Json_Decode$string)))),
	A2(
		_concourse$atc$Concourse_Job$optional,
		false,
		_elm_lang$core$Json_Decode$maybe(
			A2(_elm_lang$core$Json_Decode_ops[':='], 'trigger', _elm_lang$core$Json_Decode$bool))));
var _concourse$atc$Concourse_Job$Output = F2(
	function (a, b) {
		return {name: a, resource: b};
	});
var _concourse$atc$Concourse_Job$decodeOutput = A3(
	_elm_lang$core$Json_Decode$object2,
	_concourse$atc$Concourse_Job$Output,
	A2(_elm_lang$core$Json_Decode_ops[':='], 'name', _elm_lang$core$Json_Decode$string),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'resource', _elm_lang$core$Json_Decode$string));
var _concourse$atc$Concourse_Job$decode = F2(
	function (teamName, pipelineName) {
		return A2(
			_elm_community$json_extra$Json_Decode_Extra_ops['|:'],
			A2(
				_elm_community$json_extra$Json_Decode_Extra_ops['|:'],
				A2(
					_elm_community$json_extra$Json_Decode_Extra_ops['|:'],
					A2(
						_elm_community$json_extra$Json_Decode_Extra_ops['|:'],
						A2(
							_elm_community$json_extra$Json_Decode_Extra_ops['|:'],
							A2(
								_elm_community$json_extra$Json_Decode_Extra_ops['|:'],
								A2(
									_elm_community$json_extra$Json_Decode_Extra_ops['|:'],
									A2(
										_elm_community$json_extra$Json_Decode_Extra_ops['|:'],
										A2(
											_elm_community$json_extra$Json_Decode_Extra_ops['|:'],
											_elm_lang$core$Json_Decode$succeed(
												A2(_concourse$atc$Concourse_Job$Job, teamName, pipelineName)),
											A2(_elm_lang$core$Json_Decode_ops[':='], 'name', _elm_lang$core$Json_Decode$string)),
										A2(_elm_lang$core$Json_Decode_ops[':='], 'url', _elm_lang$core$Json_Decode$string)),
									_elm_lang$core$Json_Decode$maybe(
										A2(_elm_lang$core$Json_Decode_ops[':='], 'next_build', _concourse$atc$Concourse_Build$decode))),
								_elm_lang$core$Json_Decode$maybe(
									A2(_elm_lang$core$Json_Decode_ops[':='], 'finished_build', _concourse$atc$Concourse_Build$decode))),
							A2(
								_concourse$atc$Concourse_Job$optional,
								false,
								_elm_lang$core$Json_Decode$maybe(
									A2(_elm_lang$core$Json_Decode_ops[':='], 'paused', _elm_lang$core$Json_Decode$bool)))),
						A2(
							_concourse$atc$Concourse_Job$optional,
							false,
							_elm_lang$core$Json_Decode$maybe(
								A2(_elm_lang$core$Json_Decode_ops[':='], 'disable_manual_trigger', _elm_lang$core$Json_Decode$bool)))),
					A2(
						_concourse$atc$Concourse_Job$optional,
						_elm_lang$core$Native_List.fromArray(
							[]),
						_elm_lang$core$Json_Decode$maybe(
							A2(
								_elm_lang$core$Json_Decode_ops[':='],
								'inputs',
								_elm_lang$core$Json_Decode$list(_concourse$atc$Concourse_Job$decodeInput))))),
				A2(
					_concourse$atc$Concourse_Job$optional,
					_elm_lang$core$Native_List.fromArray(
						[]),
					_elm_lang$core$Json_Decode$maybe(
						A2(
							_elm_lang$core$Json_Decode_ops[':='],
							'outputs',
							_elm_lang$core$Json_Decode$list(_concourse$atc$Concourse_Job$decodeOutput))))),
			A2(
				_concourse$atc$Concourse_Job$optional,
				_elm_lang$core$Native_List.fromArray(
					[]),
				_elm_lang$core$Json_Decode$maybe(
					A2(
						_elm_lang$core$Json_Decode_ops[':='],
						'groups',
						_elm_lang$core$Json_Decode$list(_elm_lang$core$Json_Decode$string)))));
	});
var _concourse$atc$Concourse_Job$fetchJob = function (job) {
	return A2(
		_evancz$elm_http$Http$get,
		A2(_concourse$atc$Concourse_Job$decode, job.teamName, job.pipelineName),
		A2(
			_elm_lang$core$Basics_ops['++'],
			'/api/v1/teams/',
			A2(
				_elm_lang$core$Basics_ops['++'],
				job.teamName,
				A2(
					_elm_lang$core$Basics_ops['++'],
					'/pipelines/',
					A2(
						_elm_lang$core$Basics_ops['++'],
						job.pipelineName,
						A2(_elm_lang$core$Basics_ops['++'], '/jobs/', job.name))))));
};
var _concourse$atc$Concourse_Job$fetchJobs = function (_p1) {
	var _p2 = _p1;
	var _p4 = _p2.teamName;
	var _p3 = _p2.pipelineName;
	return A2(
		_evancz$elm_http$Http$get,
		_elm_lang$core$Json_Decode$list(
			A2(_concourse$atc$Concourse_Job$decode, _p4, _p3)),
		A2(
			_elm_lang$core$Basics_ops['++'],
			'/api/v1/teams/',
			A2(
				_elm_lang$core$Basics_ops['++'],
				_p4,
				A2(
					_elm_lang$core$Basics_ops['++'],
					'/pipelines/',
					A2(_elm_lang$core$Basics_ops['++'], _p3, '/jobs')))));
};
var _concourse$atc$Concourse_Job$PipelineLocator = F2(
	function (a, b) {
		return {teamName: a, pipelineName: b};
	});

var _concourse$atc$Build$decodeScrollEvent = A3(
	_elm_lang$core$Json_Decode$object2,
	F2(
		function (v0, v1) {
			return {ctor: '_Tuple2', _0: v0, _1: v1};
		}),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'deltaX', _elm_lang$core$Json_Decode$float),
	A2(_elm_lang$core$Json_Decode_ops[':='], 'deltaY', _elm_lang$core$Json_Decode$float));
var _concourse$atc$Build$durationTitle = F2(
	function (date, content) {
		return A2(
			_elm_lang$html$Html$div,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$title(
					A2(_mgold$elm_date_format$Date_Format$format, '%b', date))
				]),
			content);
	});
var _concourse$atc$Build$assertLeftButton = F2(
	function (action, button) {
		return _elm_lang$core$Native_Utils.eq(button, 0) ? _elm_lang$core$Result$Ok(action) : _elm_lang$core$Result$Err('placeholder error, nothing is wrong');
	});
var _concourse$atc$Build$overrideClick = function (action) {
	return A3(
		_elm_lang$html$Html_Events$onWithOptions,
		'click',
		{stopPropagation: true, preventDefault: true},
		A2(
			_elm_lang$core$Json_Decode$customDecoder,
			A2(_elm_lang$core$Json_Decode_ops[':='], 'button', _elm_lang$core$Json_Decode$int),
			_concourse$atc$Build$assertLeftButton(action)));
};
var _concourse$atc$Build$viewBuildPrepStatus = function (status) {
	var _p0 = status;
	switch (_p0.ctor) {
		case 'Unknown':
			return A2(
				_elm_lang$html$Html$i,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('fa fa-fw fa-circle-o-notch'),
						_elm_lang$html$Html_Attributes$title('thinking...')
					]),
				_elm_lang$core$Native_List.fromArray(
					[]));
		case 'Blocking':
			return A2(
				_elm_lang$html$Html$i,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('fa fa-fw fa-spin fa-circle-o-notch inactive'),
						_elm_lang$html$Html_Attributes$title('blocking')
					]),
				_elm_lang$core$Native_List.fromArray(
					[]));
		default:
			return A2(
				_elm_lang$html$Html$i,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('fa fa-fw fa-check'),
						_elm_lang$html$Html_Attributes$title('not blocking')
					]),
				_elm_lang$core$Native_List.fromArray(
					[]));
	}
};
var _concourse$atc$Build$viewDetailItem = function (_p1) {
	var _p2 = _p1;
	return A2(
		_elm_lang$html$Html$li,
		_elm_lang$core$Native_List.fromArray(
			[]),
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html$text(
				A2(
					_elm_lang$core$Basics_ops['++'],
					_p2._0,
					A2(_elm_lang$core$Basics_ops['++'], ' - ', _p2._1)))
			]));
};
var _concourse$atc$Build$viewBuildPrepDetails = function (details) {
	return A2(
		_elm_lang$html$Html$ul,
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html_Attributes$class('details')
			]),
		A2(
			_elm_lang$core$List$map,
			_concourse$atc$Build$viewDetailItem,
			_elm_lang$core$Dict$toList(details)));
};
var _concourse$atc$Build$viewBuildPrepLi = F3(
	function (text, status, details) {
		return A2(
			_elm_lang$html$Html$li,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$classList(
					_elm_lang$core$Native_List.fromArray(
						[
							{ctor: '_Tuple2', _0: 'prep-status', _1: true},
							{
							ctor: '_Tuple2',
							_0: 'inactive',
							_1: _elm_lang$core$Native_Utils.eq(status, _concourse$atc$Concourse_BuildPrep$Unknown)
						}
						]))
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(
					_elm_lang$html$Html$span,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('marker')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							_concourse$atc$Build$viewBuildPrepStatus(status)
						])),
					A2(
					_elm_lang$html$Html$span,
					_elm_lang$core$Native_List.fromArray(
						[]),
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html$text(text)
						])),
					_concourse$atc$Build$viewBuildPrepDetails(details)
				]));
	});
var _concourse$atc$Build$viewBuildPrepInput = function (_p3) {
	var _p4 = _p3;
	return A3(
		_concourse$atc$Build$viewBuildPrepLi,
		A2(_elm_lang$core$Basics_ops['++'], 'discovering any new versions of ', _p4._0),
		_p4._1,
		_elm_lang$core$Dict$empty);
};
var _concourse$atc$Build$viewBuildPrepInputs = function (inputs) {
	return A2(
		_elm_lang$core$List$map,
		_concourse$atc$Build$viewBuildPrepInput,
		_elm_lang$core$Dict$toList(inputs));
};
var _concourse$atc$Build$viewBuildPrep = function (prep) {
	var _p5 = prep;
	if (_p5.ctor === 'Just') {
		var _p6 = _p5._0;
		return A2(
			_elm_lang$html$Html$div,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$class('build-step')
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(
					_elm_lang$html$Html$div,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('header')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$html$Html$i,
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html_Attributes$class('left fa fa-fw fa-cogs')
								]),
							_elm_lang$core$Native_List.fromArray(
								[])),
							A2(
							_elm_lang$html$Html$h3,
							_elm_lang$core$Native_List.fromArray(
								[]),
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html$text('preparing build')
								]))
						])),
					A2(
					_elm_lang$html$Html$div,
					_elm_lang$core$Native_List.fromArray(
						[]),
					_elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$html$Html$ul,
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html_Attributes$class('prep-status-list')
								]),
							A2(
								_elm_lang$core$Basics_ops['++'],
								_elm_lang$core$Native_List.fromArray(
									[
										A3(_concourse$atc$Build$viewBuildPrepLi, 'checking pipeline is not paused', _p6.pausedPipeline, _elm_lang$core$Dict$empty),
										A3(_concourse$atc$Build$viewBuildPrepLi, 'checking job is not paused', _p6.pausedJob, _elm_lang$core$Dict$empty)
									]),
								A2(
									_elm_lang$core$Basics_ops['++'],
									_concourse$atc$Build$viewBuildPrepInputs(_p6.inputs),
									_elm_lang$core$Native_List.fromArray(
										[
											A3(_concourse$atc$Build$viewBuildPrepLi, 'waiting for a suitable set of input versions', _p6.inputsSatisfied, _p6.missingInputReasons),
											A3(_concourse$atc$Build$viewBuildPrepLi, 'checking max-in-flight is not reached', _p6.maxRunningBuilds, _elm_lang$core$Dict$empty)
										]))))
						]))
				]));
	} else {
		return A2(
			_elm_lang$html$Html$div,
			_elm_lang$core$Native_List.fromArray(
				[]),
			_elm_lang$core$Native_List.fromArray(
				[]));
	}
};
var _concourse$atc$Build$mmDDYY = function (d) {
	return A2(
		_elm_lang$core$Basics_ops['++'],
		A2(_mgold$elm_date_format$Date_Format$format, '%m/%d/', d),
		A2(
			_elm_lang$core$String$right,
			2,
			A2(_mgold$elm_date_format$Date_Format$format, '%Y', d)));
};
var _concourse$atc$Build$handleBuildPrepFetched = F2(
	function (buildPrep, model) {
		return {
			ctor: '_Tuple2',
			_0: _elm_lang$core$Native_Utils.update(
				model,
				{
					currentBuild: A2(
						_elm_lang$core$Maybe$map,
						function (info) {
							return _elm_lang$core$Native_Utils.update(
								info,
								{
									prep: _elm_lang$core$Maybe$Just(buildPrep)
								});
						},
						model.currentBuild)
				}),
			_1: _elm_lang$core$Platform_Cmd$none
		};
	});
var _concourse$atc$Build$handleBuildJobFetched = F2(
	function (job, model) {
		var withJobDetails = _elm_lang$core$Native_Utils.update(
			model,
			{
				job: _elm_lang$core$Maybe$Just(job)
			});
		return {ctor: '_Tuple2', _0: withJobDetails, _1: _elm_lang$core$Platform_Cmd$none};
	});
var _concourse$atc$Build$currentBuildOutput = function (cb) {
	var _p7 = cb.buildOrOutput;
	if (_p7.ctor === 'NoOutput') {
		return _elm_lang$core$Maybe$Nothing;
	} else {
		return _elm_lang$core$Maybe$Just(_p7._0);
	}
};
var _concourse$atc$Build$currentBuildBuild = function (cb) {
	var _p8 = cb.buildOrOutput;
	if (_p8.ctor === 'NoOutput') {
		return _p8._0;
	} else {
		return _p8._0.build;
	}
};
var _concourse$atc$Build$getScrollBehavior = function (model) {
	var _p10 = A2(
		_elm_lang$core$Maybe$withDefault,
		_concourse$atc$Concourse_BuildStatus$Pending,
		A2(
			_elm_lang$core$Maybe$map,
			function (_p9) {
				return function (_) {
					return _.status;
				}(
					_concourse$atc$Build$currentBuildBuild(_p9));
			},
			model.currentBuild));
	switch (_p10.ctor) {
		case 'Failed':
			return _concourse$atc$Autoscroll$ScrollUntilCancelled;
		case 'Errored':
			return _concourse$atc$Autoscroll$ScrollUntilCancelled;
		case 'Aborted':
			return _concourse$atc$Autoscroll$ScrollUntilCancelled;
		case 'Started':
			return _concourse$atc$Autoscroll$Autoscroll;
		case 'Pending':
			return _concourse$atc$Autoscroll$NoScroll;
		default:
			return _concourse$atc$Autoscroll$NoScroll;
	}
};
var _concourse$atc$Build$CurrentBuild = F2(
	function (a, b) {
		return {buildOrOutput: a, prep: b};
	});
var _concourse$atc$Build$Model = F4(
	function (a, b, c, d) {
		return {now: a, job: b, history: c, currentBuild: d};
	});
var _concourse$atc$Build$Flags = function (a) {
	return {buildId: a};
};
var _concourse$atc$Build$NoOutput = function (a) {
	return {ctor: 'NoOutput', _0: a};
};
var _concourse$atc$Build$HasOutput = function (a) {
	return {ctor: 'HasOutput', _0: a};
};
var _concourse$atc$Build$LoginRequired = {ctor: 'LoginRequired'};
var _concourse$atc$Build$StepsComplete = {ctor: 'StepsComplete'};
var _concourse$atc$Build$StepsLiveUpdating = {ctor: 'StepsLiveUpdating'};
var _concourse$atc$Build$StepsLoading = {ctor: 'StepsLoading'};
var _concourse$atc$Build$RevealCurrentBuildInHistory = {ctor: 'RevealCurrentBuildInHistory'};
var _concourse$atc$Build$BuildAborted = function (a) {
	return {ctor: 'BuildAborted', _0: a};
};
var _concourse$atc$Build$abortBuild = function (buildId) {
	return function (_p11) {
		return A2(
			_elm_lang$core$Platform_Cmd$map,
			_concourse$atc$Build$BuildAborted,
			A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _p11));
	}(
		_concourse$atc$Concourse_Build$abort(buildId));
};
var _concourse$atc$Build$ClockTick = function (a) {
	return {ctor: 'ClockTick', _0: a};
};
var _concourse$atc$Build$ScrollBuilds = function (a) {
	return {ctor: 'ScrollBuilds', _0: a};
};
var _concourse$atc$Build$BuildOutputAction = function (a) {
	return {ctor: 'BuildOutputAction', _0: a};
};
var _concourse$atc$Build$initBuildOutput = F2(
	function (build, model) {
		var _p12 = _concourse$atc$BuildOutput$init(build);
		var output = _p12._0;
		var outputCmd = _p12._1;
		return {
			ctor: '_Tuple2',
			_0: _elm_lang$core$Native_Utils.update(
				model,
				{
					currentBuild: A2(
						_elm_lang$core$Maybe$map,
						function (info) {
							return _elm_lang$core$Native_Utils.update(
								info,
								{
									buildOrOutput: _concourse$atc$Build$HasOutput(output)
								});
						},
						model.currentBuild)
				}),
			_1: A2(_elm_lang$core$Platform_Cmd$map, _concourse$atc$Build$BuildOutputAction, outputCmd)
		};
	});
var _concourse$atc$Build$viewBuildOutput = function (output) {
	var _p13 = output;
	if (_p13.ctor === 'Just') {
		return A2(
			_elm_lang$html$Html_App$map,
			_concourse$atc$Build$BuildOutputAction,
			_concourse$atc$BuildOutput$view(_p13._0));
	} else {
		return A2(
			_elm_lang$html$Html$div,
			_elm_lang$core$Native_List.fromArray(
				[]),
			_elm_lang$core$Native_List.fromArray(
				[]));
	}
};
var _concourse$atc$Build$BuildJobDetailsFetched = function (a) {
	return {ctor: 'BuildJobDetailsFetched', _0: a};
};
var _concourse$atc$Build$fetchBuildJobDetails = function (buildJob) {
	return function (_p14) {
		return A2(
			_elm_lang$core$Platform_Cmd$map,
			_concourse$atc$Build$BuildJobDetailsFetched,
			A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _p14));
	}(
		_concourse$atc$Concourse_Job$fetchJob(buildJob));
};
var _concourse$atc$Build$BuildHistoryFetched = function (a) {
	return {ctor: 'BuildHistoryFetched', _0: a};
};
var _concourse$atc$Build$fetchBuildHistory = F2(
	function (job, page) {
		return function (_p15) {
			return A2(
				_elm_lang$core$Platform_Cmd$map,
				_concourse$atc$Build$BuildHistoryFetched,
				A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _p15));
		}(
			A2(_concourse$atc$Concourse_Build$fetchJobBuilds, job, page));
	});
var _concourse$atc$Build$handleHistoryFetched = F2(
	function (history, model) {
		var withBuilds = _elm_lang$core$Native_Utils.update(
			model,
			{
				history: A2(_elm_lang$core$List$append, model.history, history.content)
			});
		var _p17 = {
			ctor: '_Tuple2',
			_0: history.pagination.nextPage,
			_1: A2(
				_elm_lang$core$Maybe$andThen,
				model.currentBuild,
				function (_p16) {
					return function (_) {
						return _.job;
					}(
						_concourse$atc$Build$currentBuildBuild(_p16));
				})
		};
		if (_p17._0.ctor === 'Nothing') {
			return {ctor: '_Tuple2', _0: withBuilds, _1: _elm_lang$core$Platform_Cmd$none};
		} else {
			if (_p17._1.ctor === 'Just') {
				return {
					ctor: '_Tuple2',
					_0: withBuilds,
					_1: _elm_lang$core$Platform_Cmd$batch(
						_elm_lang$core$Native_List.fromArray(
							[
								A2(
								_concourse$atc$Build$fetchBuildHistory,
								_p17._1._0,
								_elm_lang$core$Maybe$Just(_p17._0._0))
							]))
				};
			} else {
				return _elm_lang$core$Native_Utils.crashCase(
					'Build',
					{
						start: {line: 250, column: 5},
						end: {line: 258, column: 33}
					},
					_p17)('impossible');
			}
		}
	});
var _concourse$atc$Build$BuildPrepFetched = function (a) {
	return {ctor: 'BuildPrepFetched', _0: a};
};
var _concourse$atc$Build$fetchBuildPrep = F2(
	function (delay, buildId) {
		return function (_p19) {
			return A2(
				_elm_lang$core$Platform_Cmd$map,
				_concourse$atc$Build$BuildPrepFetched,
				A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _p19));
		}(
			A2(
				_elm_lang$core$Task$andThen,
				_elm_lang$core$Process$sleep(delay),
				_elm_lang$core$Basics$always(
					_concourse$atc$Concourse_BuildPrep$fetch(buildId))));
	});
var _concourse$atc$Build$FinishedWaitingForBuildFetch = function (a) {
	return {ctor: 'FinishedWaitingForBuildFetch', _0: a};
};
var _concourse$atc$Build$pollUntilStarted = function (buildId) {
	return A3(
		_elm_lang$core$Task$perform,
		_elm_lang$core$Basics$identity,
		_elm_lang$core$Basics$identity,
		A2(
			_elm_lang$core$Task$andThen,
			_elm_lang$core$Process$sleep(_elm_lang$core$Time$second),
			_elm_lang$core$Basics$always(
				_elm_lang$core$Task$succeed(
					_concourse$atc$Build$FinishedWaitingForBuildFetch(buildId)))));
};
var _concourse$atc$Build$handleBuildFetched = F2(
	function (build, model) {
		var updateBuildURL = _elm_lang$navigation$Navigation$newUrl(
			_concourse$atc$Concourse_Build$url(build));
		var fetchJobAndHistory = function () {
			var _p20 = {ctor: '_Tuple2', _0: model.job, _1: build.job};
			if (((_p20.ctor === '_Tuple2') && (_p20._0.ctor === 'Nothing')) && (_p20._1.ctor === 'Just')) {
				var _p21 = _p20._1._0;
				return _elm_lang$core$Platform_Cmd$batch(
					_elm_lang$core$Native_List.fromArray(
						[
							_concourse$atc$Build$fetchBuildJobDetails(_p21),
							A2(_concourse$atc$Build$fetchBuildHistory, _p21, _elm_lang$core$Maybe$Nothing)
						]));
			} else {
				return _elm_lang$core$Platform_Cmd$none;
			}
		}();
		var currentBuild = function () {
			var _p22 = model.currentBuild;
			if (_p22.ctor === 'Nothing') {
				return {
					buildOrOutput: _concourse$atc$Build$NoOutput(build),
					prep: _elm_lang$core$Maybe$Nothing
				};
			} else {
				return _elm_lang$core$Native_Utils.update(
					_p22._0,
					{
						buildOrOutput: _concourse$atc$Build$NoOutput(build)
					});
			}
		}();
		var withBuild = _elm_lang$core$Native_Utils.update(
			model,
			{
				currentBuild: _elm_lang$core$Maybe$Just(currentBuild)
			});
		var _p23 = function () {
			if (_elm_lang$core$Native_Utils.eq(build.status, _concourse$atc$Concourse_BuildStatus$Pending)) {
				return {
					ctor: '_Tuple2',
					_0: withBuild,
					_1: _concourse$atc$Build$pollUntilStarted(build.id)
				};
			} else {
				if (_elm_lang$core$Native_Utils.eq(build.reapTime, _elm_lang$core$Maybe$Nothing)) {
					var _p24 = A2(
						_elm_lang$core$Maybe$andThen,
						model.currentBuild,
						function (_) {
							return _.prep;
						});
					if (_p24.ctor === 'Nothing') {
						return A2(_concourse$atc$Build$initBuildOutput, build, withBuild);
					} else {
						var _p25 = A2(_concourse$atc$Build$initBuildOutput, build, withBuild);
						var newModel = _p25._0;
						var cmd = _p25._1;
						return {
							ctor: '_Tuple2',
							_0: newModel,
							_1: _elm_lang$core$Platform_Cmd$batch(
								_elm_lang$core$Native_List.fromArray(
									[
										cmd,
										A2(_concourse$atc$Build$fetchBuildPrep, _elm_lang$core$Time$second, build.id)
									]))
						};
					}
				} else {
					return {ctor: '_Tuple2', _0: withBuild, _1: _elm_lang$core$Platform_Cmd$none};
				}
			}
		}();
		var newModel = _p23._0;
		var cmd = _p23._1;
		return {
			ctor: '_Tuple2',
			_0: newModel,
			_1: _elm_lang$core$Platform_Cmd$batch(
				_elm_lang$core$Native_List.fromArray(
					[cmd, updateBuildURL, fetchJobAndHistory]))
		};
	});
var _concourse$atc$Build$BuildFetched = function (a) {
	return {ctor: 'BuildFetched', _0: a};
};
var _concourse$atc$Build$fetchBuild = function (buildId) {
	return function (_p26) {
		return A2(
			_elm_lang$core$Platform_Cmd$map,
			_concourse$atc$Build$BuildFetched,
			A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _p26));
	}(
		_concourse$atc$Concourse_Build$fetch(buildId));
};
var _concourse$atc$Build$fetchIfPending = F2(
	function (buildId, model) {
		var _p28 = A2(
			_elm_lang$core$Maybe$withDefault,
			_concourse$atc$Concourse_BuildStatus$Pending,
			A2(
				_elm_lang$core$Maybe$map,
				function (_p27) {
					return function (_) {
						return _.status;
					}(
						_concourse$atc$Build$currentBuildBuild(_p27));
				},
				model.currentBuild));
		if (_p28.ctor === 'Pending') {
			return _elm_lang$core$Platform_Cmd$batch(
				_elm_lang$core$Native_List.fromArray(
					[
						_concourse$atc$Build$fetchBuild(buildId),
						A2(_concourse$atc$Build$fetchBuildPrep, 0, buildId)
					]));
		} else {
			return _elm_lang$core$Platform_Cmd$none;
		}
	});
var _concourse$atc$Build$AbortBuild = function (a) {
	return {ctor: 'AbortBuild', _0: a};
};
var _concourse$atc$Build$FetchBuild = function (a) {
	return {ctor: 'FetchBuild', _0: a};
};
var _concourse$atc$Build$viewHistoryItem = F2(
	function (currentBuild, build) {
		return A2(
			_elm_lang$html$Html$li,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$classList(
					_elm_lang$core$Native_List.fromArray(
						[
							{
							ctor: '_Tuple2',
							_0: _elm_lang$core$Native_Utils.eq(build.name, currentBuild.name) ? _concourse$atc$Concourse_BuildStatus$show(currentBuild.status) : _concourse$atc$Concourse_BuildStatus$show(build.status),
							_1: true
						},
							{
							ctor: '_Tuple2',
							_0: 'current',
							_1: _elm_lang$core$Native_Utils.eq(build.name, currentBuild.name)
						}
						]))
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(
					_elm_lang$html$Html$a,
					_elm_lang$core$Native_List.fromArray(
						[
							_concourse$atc$Build$overrideClick(
							_concourse$atc$Build$FetchBuild(build.id)),
							_elm_lang$html$Html_Attributes$href(
							_concourse$atc$Concourse_Build$url(build))
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html$text(build.name)
						]))
				]));
	});
var _concourse$atc$Build$viewHistory = F2(
	function (currentBuild, builds) {
		return A2(
			_elm_lang$html$Html$ul,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$id('builds')
				]),
			A2(
				_elm_lang$core$List$map,
				_concourse$atc$Build$viewHistoryItem(currentBuild),
				builds));
	});
var _concourse$atc$Build$lazyViewHistory = F2(
	function (currentBuild, builds) {
		return A3(_elm_lang$html$Html_Lazy$lazy2, _concourse$atc$Build$viewHistory, currentBuild, builds);
	});
var _concourse$atc$Build$viewBuildHeader = F2(
	function (build, _p29) {
		var _p30 = _p29;
		var _p35 = _p30.job;
		var buildTitle = function () {
			var _p31 = build.job;
			if (_p31.ctor === 'Just') {
				var _p32 = _p31._0.name;
				return A2(
					_elm_lang$html$Html$a,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$href(
							A2(
								_elm_lang$core$Basics_ops['++'],
								'/teams/',
								A2(
									_elm_lang$core$Basics_ops['++'],
									_p31._0.teamName,
									A2(
										_elm_lang$core$Basics_ops['++'],
										'/pipelines/',
										A2(
											_elm_lang$core$Basics_ops['++'],
											_p31._0.pipelineName,
											A2(_elm_lang$core$Basics_ops['++'], '/jobs/', _p32))))))
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html$text(
							A2(
								_elm_lang$core$Basics_ops['++'],
								_p32,
								A2(_elm_lang$core$Basics_ops['++'], ' #', build.name)))
						]));
			} else {
				return _elm_lang$html$Html$text(
					A2(
						_elm_lang$core$Basics_ops['++'],
						'build #',
						_elm_lang$core$Basics$toString(build.id)));
			}
		}();
		var abortButton = _concourse$atc$Concourse_BuildStatus$isRunning(build.status) ? A2(
			_elm_lang$html$Html$span,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$class('build-action build-action-abort fr'),
					_elm_lang$html$Html_Events$onClick(
					_concourse$atc$Build$AbortBuild(build.id)),
					A2(_elm_lang$html$Html_Attributes$attribute, 'aria-label', 'Abort Build')
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(
					_elm_lang$html$Html$i,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('fa fa-times-circle')
						]),
					_elm_lang$core$Native_List.fromArray(
						[]))
				])) : A2(
			_elm_lang$html$Html$span,
			_elm_lang$core$Native_List.fromArray(
				[]),
			_elm_lang$core$Native_List.fromArray(
				[]));
		var triggerButton = function () {
			var _p33 = _p35;
			if (_p33.ctor === 'Just') {
				var buttonDisabled = function () {
					var _p34 = _p35;
					if (_p34.ctor === 'Nothing') {
						return true;
					} else {
						return _p34._0.disableManualTrigger;
					}
				}();
				var actionUrl = A2(
					_elm_lang$core$Basics_ops['++'],
					'/teams/',
					A2(
						_elm_lang$core$Basics_ops['++'],
						_p33._0.teamName,
						A2(
							_elm_lang$core$Basics_ops['++'],
							'/pipelines/',
							A2(
								_elm_lang$core$Basics_ops['++'],
								_p33._0.pipelineName,
								A2(
									_elm_lang$core$Basics_ops['++'],
									'/jobs/',
									A2(_elm_lang$core$Basics_ops['++'], _p33._0.name, '/builds'))))));
				return A2(
					_elm_lang$html$Html$form,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('trigger-build'),
							_elm_lang$html$Html_Attributes$method('post'),
							_elm_lang$html$Html_Attributes$action(actionUrl)
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$html$Html$button,
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html_Attributes$class('build-action fr'),
									_elm_lang$html$Html_Attributes$disabled(buttonDisabled),
									A2(_elm_lang$html$Html_Attributes$attribute, 'aria-label', 'Trigger Build')
								]),
							_elm_lang$core$Native_List.fromArray(
								[
									A2(
									_elm_lang$html$Html$i,
									_elm_lang$core$Native_List.fromArray(
										[
											_elm_lang$html$Html_Attributes$class('fa fa-plus-circle')
										]),
									_elm_lang$core$Native_List.fromArray(
										[]))
								]))
						]));
			} else {
				return A2(
					_elm_lang$html$Html$div,
					_elm_lang$core$Native_List.fromArray(
						[]),
					_elm_lang$core$Native_List.fromArray(
						[]));
			}
		}();
		return A2(
			_elm_lang$html$Html$div,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$class('fixed-header')
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(
					_elm_lang$html$Html$div,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class(
							A2(
								_elm_lang$core$Basics_ops['++'],
								'build-header ',
								_concourse$atc$Concourse_BuildStatus$show(build.status)))
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$html$Html$div,
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html_Attributes$class('build-actions fr')
								]),
							_elm_lang$core$Native_List.fromArray(
								[triggerButton, abortButton])),
							A2(
							_elm_lang$html$Html$h1,
							_elm_lang$core$Native_List.fromArray(
								[]),
							_elm_lang$core$Native_List.fromArray(
								[buildTitle])),
							A2(_concourse$atc$BuildDuration$view, build.duration, _p30.now)
						])),
					A2(
					_elm_lang$html$Html$div,
					_elm_lang$core$Native_List.fromArray(
						[
							A3(
							_elm_lang$html$Html_Events$onWithOptions,
							'mousewheel',
							{stopPropagation: true, preventDefault: true},
							A2(_elm_lang$core$Json_Decode$map, _concourse$atc$Build$ScrollBuilds, _concourse$atc$Build$decodeScrollEvent))
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							A2(_concourse$atc$Build$lazyViewHistory, build, _p30.history)
						]))
				]));
	});
var _concourse$atc$Build$view = function (model) {
	var _p36 = model.currentBuild;
	if (_p36.ctor === 'Just') {
		var _p40 = _p36._0;
		return A2(
			_elm_lang$html$Html$div,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$class('with-fixed-header')
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(
					_concourse$atc$Build$viewBuildHeader,
					_concourse$atc$Build$currentBuildBuild(_p40),
					model),
					A2(
					_elm_lang$html$Html$div,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('scrollable-body')
						]),
					A2(
						_elm_lang$core$Basics_ops['++'],
						_elm_lang$core$Native_List.fromArray(
							[
								_concourse$atc$Build$viewBuildPrep(_p40.prep),
								A2(
								_elm_lang$html$Html_Lazy$lazy,
								_concourse$atc$Build$viewBuildOutput,
								_concourse$atc$Build$currentBuildOutput(_p40))
							]),
						function () {
							var build = _concourse$atc$Build$currentBuildBuild(_p40);
							var maybeBirthDate = _elm_lang$core$Maybe$oneOf(
								_elm_lang$core$Native_List.fromArray(
									[build.duration.startedAt, build.duration.finishedAt]));
							var _p37 = {ctor: '_Tuple2', _0: maybeBirthDate, _1: build.reapTime};
							if (((_p37.ctor === '_Tuple2') && (_p37._0.ctor === 'Just')) && (_p37._1.ctor === 'Just')) {
								return _elm_lang$core$Native_List.fromArray(
									[
										A2(
										_elm_lang$html$Html$div,
										_elm_lang$core$Native_List.fromArray(
											[
												_elm_lang$html$Html_Attributes$class('tombstone')
											]),
										_elm_lang$core$Native_List.fromArray(
											[
												A2(
												_elm_lang$html$Html$div,
												_elm_lang$core$Native_List.fromArray(
													[
														_elm_lang$html$Html_Attributes$class('heading')
													]),
												_elm_lang$core$Native_List.fromArray(
													[
														_elm_lang$html$Html$text('RIP')
													])),
												A2(
												_elm_lang$html$Html$div,
												_elm_lang$core$Native_List.fromArray(
													[
														_elm_lang$html$Html_Attributes$class('job-name')
													]),
												_elm_lang$core$Native_List.fromArray(
													[
														_elm_lang$html$Html$text(
														A2(
															_elm_lang$core$Maybe$withDefault,
															'one-off build',
															A2(
																_elm_lang$core$Maybe$map,
																function (_) {
																	return _.name;
																},
																build.job)))
													])),
												A2(
												_elm_lang$html$Html$div,
												_elm_lang$core$Native_List.fromArray(
													[
														_elm_lang$html$Html_Attributes$class('build-name')
													]),
												_elm_lang$core$Native_List.fromArray(
													[
														_elm_lang$html$Html$text(
														A2(
															_elm_lang$core$Basics_ops['++'],
															'build #',
															function () {
																var _p38 = build.job;
																if (_p38.ctor === 'Nothing') {
																	return _elm_lang$core$Basics$toString(build.id);
																} else {
																	return build.name;
																}
															}()))
													])),
												A2(
												_elm_lang$html$Html$div,
												_elm_lang$core$Native_List.fromArray(
													[
														_elm_lang$html$Html_Attributes$class('date')
													]),
												_elm_lang$core$Native_List.fromArray(
													[
														_elm_lang$html$Html$text(
														A2(
															_elm_lang$core$Basics_ops['++'],
															_concourse$atc$Build$mmDDYY(_p37._0._0),
															A2(
																_elm_lang$core$Basics_ops['++'],
																'-',
																_concourse$atc$Build$mmDDYY(_p37._1._0))))
													])),
												A2(
												_elm_lang$html$Html$div,
												_elm_lang$core$Native_List.fromArray(
													[
														_elm_lang$html$Html_Attributes$class('epitaph')
													]),
												_elm_lang$core$Native_List.fromArray(
													[
														_elm_lang$html$Html$text(
														function () {
															var _p39 = build.status;
															switch (_p39.ctor) {
																case 'Succeeded':
																	return 'It passed, and now it has passed on.';
																case 'Failed':
																	return 'It failed, and now has been forgotten.';
																case 'Errored':
																	return 'It errored, but has found forgiveness.';
																case 'Aborted':
																	return 'It was never given a chance.';
																default:
																	return 'I\'m not dead yet.';
															}
														}())
													]))
											])),
										A2(
										_elm_lang$html$Html$div,
										_elm_lang$core$Native_List.fromArray(
											[
												_elm_lang$html$Html_Attributes$class('explanation')
											]),
										_elm_lang$core$Native_List.fromArray(
											[
												_elm_lang$html$Html$text('This log has been '),
												A2(
												_elm_lang$html$Html$a,
												_elm_lang$core$Native_List.fromArray(
													[
														_elm_lang$html$Html_Attributes$href('http://concourse.ci/configuring-jobs.html#build_logs_to_retain')
													]),
												_elm_lang$core$Native_List.fromArray(
													[
														_elm_lang$html$Html$text('reaped.')
													]))
											]))
									]);
							} else {
								return _elm_lang$core$Native_List.fromArray(
									[]);
							}
						}()))
				]));
	} else {
		return _concourse$atc$LoadingIndicator$view;
	}
};
var _concourse$atc$Build$Noop = {ctor: 'Noop'};
var _concourse$atc$Build$scrollBuilds = function (delta) {
	return function (_p41) {
		return A2(
			_elm_lang$core$Platform_Cmd$map,
			_elm_lang$core$Basics$always(_concourse$atc$Build$Noop),
			A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _p41));
	}(
		A2(_concourse$atc$Scroll$scroll, 'builds', delta));
};
var _concourse$atc$Build$scrollToCurrentBuildInHistory = function (_p42) {
	return A2(
		_elm_lang$core$Platform_Cmd$map,
		_elm_lang$core$Basics$always(_concourse$atc$Build$Noop),
		A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _p42));
}(
	_concourse$atc$Scroll$scrollIntoView('#builds .current'));
var _concourse$atc$Build$redirectToLogin = function (model) {
	return function (_p43) {
		return A2(
			_elm_lang$core$Platform_Cmd$map,
			_elm_lang$core$Basics$always(_concourse$atc$Build$Noop),
			A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _p43));
	}(
		_concourse$atc$Redirect$to('/login'));
};
var _concourse$atc$Build$update = F2(
	function (action, model) {
		var _p44 = action;
		switch (_p44.ctor) {
			case 'Noop':
				return {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none};
			case 'FetchBuild':
				return {
					ctor: '_Tuple2',
					_0: model,
					_1: _concourse$atc$Build$fetchBuild(_p44._0)
				};
			case 'FinishedWaitingForBuildFetch':
				return {
					ctor: '_Tuple2',
					_0: model,
					_1: A2(_concourse$atc$Build$fetchIfPending, _p44._0, model)
				};
			case 'BuildFetched':
				if (_p44._0.ctor === 'Ok') {
					return A2(_concourse$atc$Build$handleBuildFetched, _p44._0._0, model);
				} else {
					return A2(
						_elm_lang$core$Debug$log,
						A2(
							_elm_lang$core$Basics_ops['++'],
							'failed to fetch build: ',
							_elm_lang$core$Basics$toString(_p44._0._0)),
						{ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none});
				}
			case 'AbortBuild':
				return {
					ctor: '_Tuple2',
					_0: model,
					_1: _concourse$atc$Build$abortBuild(_p44._0)
				};
			case 'BuildAborted':
				if (_p44._0.ctor === 'Ok') {
					return {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none};
				} else {
					if ((_p44._0._0.ctor === 'BadResponse') && (_p44._0._0._0 === 401)) {
						return {
							ctor: '_Tuple2',
							_0: model,
							_1: _concourse$atc$Build$redirectToLogin(model)
						};
					} else {
						return A2(
							_elm_lang$core$Debug$log,
							A2(
								_elm_lang$core$Basics_ops['++'],
								'failed to abort build: ',
								_elm_lang$core$Basics$toString(_p44._0._0)),
							{ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none});
					}
				}
			case 'BuildPrepFetched':
				if (_p44._0.ctor === 'Ok') {
					return A2(_concourse$atc$Build$handleBuildPrepFetched, _p44._0._0, model);
				} else {
					return A2(
						_elm_lang$core$Debug$log,
						A2(
							_elm_lang$core$Basics_ops['++'],
							'failed to fetch build preparation: ',
							_elm_lang$core$Basics$toString(_p44._0._0)),
						{ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none});
				}
			case 'BuildOutputAction':
				var _p45 = {
					ctor: '_Tuple2',
					_0: model.currentBuild,
					_1: A2(
						_elm_lang$core$Maybe$map,
						function (_) {
							return _.buildOrOutput;
						},
						model.currentBuild)
				};
				if ((((_p45.ctor === '_Tuple2') && (_p45._0.ctor === 'Just')) && (_p45._1.ctor === 'Just')) && (_p45._1._0.ctor === 'HasOutput')) {
					var _p46 = A2(_concourse$atc$BuildOutput$update, _p44._0, _p45._1._0._0);
					var newOutput = _p46._0;
					var cmd = _p46._1;
					return {
						ctor: '_Tuple2',
						_0: _elm_lang$core$Native_Utils.update(
							model,
							{
								currentBuild: _elm_lang$core$Maybe$Just(
									_elm_lang$core$Native_Utils.update(
										_p45._0._0,
										{
											buildOrOutput: _concourse$atc$Build$HasOutput(newOutput)
										}))
							}),
						_1: A2(_elm_lang$core$Platform_Cmd$map, _concourse$atc$Build$BuildOutputAction, cmd)
					};
				} else {
					return _elm_lang$core$Native_Utils.crashCase(
						'Build',
						{
							start: {line: 138, column: 7},
							end: {line: 149, column: 77}
						},
						_p45)('impossible (received action for missing BuildOutput)');
				}
			case 'BuildHistoryFetched':
				if (_p44._0.ctor === 'Err') {
					return A2(
						_elm_lang$core$Debug$log,
						A2(
							_elm_lang$core$Basics_ops['++'],
							'failed to fetch build history: ',
							_elm_lang$core$Basics$toString(_p44._0._0)),
						{ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none});
				} else {
					return A2(_concourse$atc$Build$handleHistoryFetched, _p44._0._0, model);
				}
			case 'BuildJobDetailsFetched':
				if (_p44._0.ctor === 'Ok') {
					return A2(_concourse$atc$Build$handleBuildJobFetched, _p44._0._0, model);
				} else {
					return A2(
						_elm_lang$core$Debug$log,
						A2(
							_elm_lang$core$Basics_ops['++'],
							'failed to fetch build job details: ',
							_elm_lang$core$Basics$toString(_p44._0._0)),
						{ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none});
				}
			case 'RevealCurrentBuildInHistory':
				return {ctor: '_Tuple2', _0: model, _1: _concourse$atc$Build$scrollToCurrentBuildInHistory};
			case 'ScrollBuilds':
				if (_p44._0._0 === 0) {
					return {
						ctor: '_Tuple2',
						_0: model,
						_1: _concourse$atc$Build$scrollBuilds(_p44._0._1)
					};
				} else {
					return {
						ctor: '_Tuple2',
						_0: model,
						_1: _concourse$atc$Build$scrollBuilds(0 - _p44._0._0)
					};
				}
			default:
				return {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Native_Utils.update(
						model,
						{now: _p44._0}),
					_1: _elm_lang$core$Platform_Cmd$none
				};
		}
	});
var _concourse$atc$Build$init = function (flags) {
	var model = {
		now: 0,
		job: _elm_lang$core$Maybe$Nothing,
		history: _elm_lang$core$Native_List.fromArray(
			[]),
		currentBuild: _elm_lang$core$Maybe$Nothing
	};
	return A2(
		_concourse$atc$Build$update,
		_concourse$atc$Build$FetchBuild(flags.buildId),
		model);
};

var _concourse$atc$BuildPage$main = {
	main: _elm_lang$html$Html_App$programWithFlags(
		{
			init: function (_p0) {
				return A2(
					_concourse$atc$Autoscroll$init,
					_concourse$atc$Build$getScrollBehavior,
					_concourse$atc$Build$init(_p0));
			},
			update: _concourse$atc$Autoscroll$update(_concourse$atc$Build$update),
			view: _concourse$atc$Autoscroll$view(_concourse$atc$Build$view),
			subscriptions: function () {
				var pushDown = A2(
					_elm_lang$core$Time$every,
					100 * _elm_lang$core$Time$millisecond,
					_elm_lang$core$Basics$always(_concourse$atc$Autoscroll$ScrollDown));
				var scrolledUp = _concourse$atc$Scroll$fromBottom(_concourse$atc$Autoscroll$FromBottom);
				var tick = A2(
					_elm_lang$core$Time$every,
					_elm_lang$core$Time$second,
					function (_p1) {
						return _concourse$atc$Autoscroll$SubAction(
							_concourse$atc$Build$ClockTick(_p1));
					});
				return function (model) {
					return _elm_lang$core$Platform_Sub$batch(
						_elm_lang$core$Native_List.fromArray(
							[
								tick,
								scrolledUp,
								pushDown,
								function () {
								var _p2 = A2(_elm_lang$core$Maybe$andThen, model.subModel.currentBuild, _concourse$atc$Build$currentBuildOutput);
								if (_p2.ctor === 'Nothing') {
									return _elm_lang$core$Platform_Sub$none;
								} else {
									return A2(
										_elm_lang$core$Platform_Sub$map,
										function (_p3) {
											return _concourse$atc$Autoscroll$SubAction(
												_concourse$atc$Build$BuildOutputAction(_p3));
										},
										_p2._0.events);
								}
							}()
							]));
				};
			}()
		}),
	flags: A2(
		_elm_lang$core$Json_Decode$andThen,
		A2(_elm_lang$core$Json_Decode_ops[':='], 'buildId', _elm_lang$core$Json_Decode$int),
		function (buildId) {
			return _elm_lang$core$Json_Decode$succeed(
				{buildId: buildId});
		})
};

var _concourse$atc$Concourse_AuthMethod$OAuthAuthMethod = F2(
	function (a, b) {
		return {displayName: a, authURL: b};
	});
var _concourse$atc$Concourse_AuthMethod$OAuthMethod = function (a) {
	return {ctor: 'OAuthMethod', _0: a};
};
var _concourse$atc$Concourse_AuthMethod$BasicMethod = {ctor: 'BasicMethod'};
var _concourse$atc$Concourse_AuthMethod$authMethodFromTuple = function (tuple) {
	var _p0 = tuple;
	_v0_3:
	do {
		if (_p0.ctor === '_Tuple3') {
			switch (_p0._0) {
				case 'basic':
					return _elm_lang$core$Result$Ok(_concourse$atc$Concourse_AuthMethod$BasicMethod);
				case 'oauth':
					if ((_p0._1.ctor === 'Just') && (_p0._2.ctor === 'Just')) {
						return _elm_lang$core$Result$Ok(
							_concourse$atc$Concourse_AuthMethod$OAuthMethod(
								{displayName: _p0._1._0, authURL: _p0._2._0}));
					} else {
						return _elm_lang$core$Result$Err('missing fields in oauth auth method');
					}
				default:
					break _v0_3;
			}
		} else {
			break _v0_3;
		}
	} while(false);
	return _elm_lang$core$Result$Err('unknown value for auth method type');
};
var _concourse$atc$Concourse_AuthMethod$decodeAuthMethods = _elm_lang$core$Json_Decode$list(
	A2(
		_elm_lang$core$Json_Decode$customDecoder,
		A4(
			_elm_lang$core$Json_Decode$object3,
			F3(
				function (v0, v1, v2) {
					return {ctor: '_Tuple3', _0: v0, _1: v1, _2: v2};
				}),
			A2(_elm_lang$core$Json_Decode_ops[':='], 'type', _elm_lang$core$Json_Decode$string),
			_elm_lang$core$Json_Decode$maybe(
				A2(_elm_lang$core$Json_Decode_ops[':='], 'display_name', _elm_lang$core$Json_Decode$string)),
			_elm_lang$core$Json_Decode$maybe(
				A2(_elm_lang$core$Json_Decode_ops[':='], 'auth_url', _elm_lang$core$Json_Decode$string))),
		_concourse$atc$Concourse_AuthMethod$authMethodFromTuple));
var _concourse$atc$Concourse_AuthMethod$fetchAuthMethods = function (teamName) {
	return A2(
		_evancz$elm_http$Http$get,
		_concourse$atc$Concourse_AuthMethod$decodeAuthMethods,
		A2(
			_elm_lang$core$Basics_ops['++'],
			'/api/v1/teams/',
			A2(_elm_lang$core$Basics_ops['++'], teamName, '/auth/methods')));
};

var _concourse$atc$Concourse_Team$Team = function (a) {
	return {name: a};
};
var _concourse$atc$Concourse_Team$decodeTeams = _elm_lang$core$Json_Decode$list(
	A2(
		_elm_lang$core$Json_Decode$object1,
		_concourse$atc$Concourse_Team$Team,
		A2(_elm_lang$core$Json_Decode_ops[':='], 'name', _elm_lang$core$Json_Decode$string)));
var _concourse$atc$Concourse_Team$fetchTeams = A2(_evancz$elm_http$Http$get, _concourse$atc$Concourse_Team$decodeTeams, '/api/v1/teams');

var _imeckler$queue$Queue_Internal$Queue = F2(
	function (a, b) {
		return {ctor: 'Queue', _0: a, _1: b};
	});

var _imeckler$queue$Queue$toList = function (_p0) {
	var _p1 = _p0;
	return A2(
		_elm_lang$core$Basics_ops['++'],
		_p1._0,
		_elm_lang$core$List$reverse(_p1._1));
};
var _imeckler$queue$Queue$map = F2(
	function (g, _p2) {
		var _p3 = _p2;
		return A2(
			_imeckler$queue$Queue_Internal$Queue,
			A2(_elm_lang$core$List$map, g, _p3._0),
			A2(_elm_lang$core$List$map, g, _p3._1));
	});
var _imeckler$queue$Queue$length = function (_p4) {
	var _p5 = _p4;
	return _elm_lang$core$List$length(_p5._0) + _elm_lang$core$List$length(_p5._1);
};
var _imeckler$queue$Queue$isEmpty = function (q) {
	var _p6 = q;
	if ((_p6._0.ctor === '[]') && (_p6._1.ctor === '[]')) {
		return true;
	} else {
		return false;
	}
};
var _imeckler$queue$Queue$pop = function (_p7) {
	var _p8 = _p7;
	var _p13 = _p8._1;
	var _p9 = _p8._0;
	if (_p9.ctor === '::') {
		return _elm_lang$core$Maybe$Just(
			{
				ctor: '_Tuple2',
				_0: _p9._0,
				_1: A2(_imeckler$queue$Queue_Internal$Queue, _p9._1, _p13)
			});
	} else {
		var _p10 = _p13;
		if (_p10.ctor === '[]') {
			return _elm_lang$core$Maybe$Nothing;
		} else {
			var _p11 = _elm_lang$core$List$reverse(_p13);
			if (_p11.ctor === '::') {
				return _elm_lang$core$Maybe$Just(
					{
						ctor: '_Tuple2',
						_0: _p11._0,
						_1: A2(
							_imeckler$queue$Queue_Internal$Queue,
							_p11._1,
							_elm_lang$core$Native_List.fromArray(
								[]))
					});
			} else {
				return _elm_lang$core$Native_Utils.crashCase(
					'Queue',
					{
						start: {line: 57, column: 11},
						end: {line: 59, column: 44}
					},
					_p11)('Impossible!');
			}
		}
	}
};
var _imeckler$queue$Queue$push = F2(
	function (x, _p14) {
		var _p15 = _p14;
		return A2(
			_imeckler$queue$Queue_Internal$Queue,
			_p15._0,
			A2(_elm_lang$core$List_ops['::'], x, _p15._1));
	});
var _imeckler$queue$Queue$empty = A2(
	_imeckler$queue$Queue_Internal$Queue,
	_elm_lang$core$Native_List.fromArray(
		[]),
	_elm_lang$core$Native_List.fromArray(
		[]));

var _elm_community$graph$Graph_Tree$pushMany = F2(
	function (vals, queue) {
		return A3(_elm_lang$core$List$foldl, _imeckler$queue$Queue$push, queue, vals);
	});
var _elm_community$graph$Graph_Tree$listForTraversal = F2(
	function (traversal, tree) {
		var acc = _elm_lang$core$Basics$identity;
		var f = F3(
			function (label, children, rest) {
				return function (_p0) {
					return rest(
						A2(
							F2(
								function (x, y) {
									return A2(_elm_lang$core$List_ops['::'], x, y);
								}),
							label,
							_p0));
				};
			});
		return A4(
			traversal,
			f,
			acc,
			tree,
			_elm_lang$core$Native_List.fromArray(
				[]));
	});
var _elm_community$graph$Graph_Tree$size = function (tree) {
	var _p1 = tree;
	return _p1._0;
};
var _elm_community$graph$Graph_Tree$root = function (tree) {
	var _p2 = tree;
	return _p2._1;
};
var _elm_community$graph$Graph_Tree$height = function (tree) {
	var go = F2(
		function (h, t) {
			var _p3 = _elm_community$graph$Graph_Tree$root(t);
			if (_p3.ctor === 'Just') {
				return A3(
					_elm_lang$core$List$foldl,
					function (_p4) {
						return _elm_lang$core$Basics$max(
							A2(go, h + 1, _p4));
					},
					h + 1,
					_p3._0._1);
			} else {
				return h;
			}
		});
	return A2(go, 0, tree);
};
var _elm_community$graph$Graph_Tree$levelOrder = F3(
	function (visit, acc, tree) {
		var go = F2(
			function (acc, toVisit) {
				go:
				while (true) {
					var _p5 = _imeckler$queue$Queue$pop(toVisit);
					if (_p5.ctor === 'Nothing') {
						return acc;
					} else {
						var _p8 = _p5._0._1;
						var _p6 = _elm_community$graph$Graph_Tree$root(_p5._0._0);
						if (_p6.ctor === 'Nothing') {
							var _v5 = acc,
								_v6 = _p8;
							acc = _v5;
							toVisit = _v6;
							continue go;
						} else {
							var _p7 = _p6._0._1;
							var _v7 = A3(visit, _p6._0._0, _p7, acc),
								_v8 = A2(_elm_community$graph$Graph_Tree$pushMany, _p7, _p8);
							acc = _v7;
							toVisit = _v8;
							continue go;
						}
					}
				}
			});
		return A2(
			go,
			acc,
			A2(_imeckler$queue$Queue$push, tree, _imeckler$queue$Queue$empty));
	});
var _elm_community$graph$Graph_Tree$levelOrderList = _elm_community$graph$Graph_Tree$listForTraversal(_elm_community$graph$Graph_Tree$levelOrder);
var _elm_community$graph$Graph_Tree$postOrder = F3(
	function (visit, acc, tree) {
		var folder = _elm_lang$core$Basics$flip(
			_elm_community$graph$Graph_Tree$postOrder(visit));
		var _p9 = _elm_community$graph$Graph_Tree$root(tree);
		if (_p9.ctor === 'Nothing') {
			return acc;
		} else {
			var _p10 = _p9._0._1;
			return A3(
				visit,
				_p9._0._0,
				_p10,
				A3(_elm_lang$core$List$foldl, folder, acc, _p10));
		}
	});
var _elm_community$graph$Graph_Tree$postOrderList = _elm_community$graph$Graph_Tree$listForTraversal(_elm_community$graph$Graph_Tree$postOrder);
var _elm_community$graph$Graph_Tree$preOrder = F3(
	function (visit, acc, tree) {
		var folder = _elm_lang$core$Basics$flip(
			_elm_community$graph$Graph_Tree$preOrder(visit));
		var _p11 = _elm_community$graph$Graph_Tree$root(tree);
		if (_p11.ctor === 'Nothing') {
			return acc;
		} else {
			var _p12 = _p11._0._1;
			return A3(
				_elm_lang$core$List$foldl,
				folder,
				A3(visit, _p11._0._0, _p12, acc),
				_p12);
		}
	});
var _elm_community$graph$Graph_Tree$preOrderList = _elm_community$graph$Graph_Tree$listForTraversal(_elm_community$graph$Graph_Tree$preOrder);
var _elm_community$graph$Graph_Tree$MkTree = F2(
	function (a, b) {
		return {ctor: 'MkTree', _0: a, _1: b};
	});
var _elm_community$graph$Graph_Tree$empty = A2(_elm_community$graph$Graph_Tree$MkTree, 0, _elm_lang$core$Maybe$Nothing);
var _elm_community$graph$Graph_Tree$isEmpty = function (tree) {
	return _elm_lang$core$Native_Utils.eq(tree, _elm_community$graph$Graph_Tree$empty);
};
var _elm_community$graph$Graph_Tree$inner = F2(
	function (label, children) {
		var children$ = A2(
			_elm_lang$core$List$filter,
			function (_p13) {
				return _elm_lang$core$Basics$not(
					_elm_community$graph$Graph_Tree$isEmpty(_p13));
			},
			children);
		var size$ = A3(
			_elm_lang$core$List$foldl,
			function (_p14) {
				return F2(
					function (x, y) {
						return x + y;
					})(
					_elm_community$graph$Graph_Tree$size(_p14));
			},
			1,
			children$);
		return A2(
			_elm_community$graph$Graph_Tree$MkTree,
			size$,
			_elm_lang$core$Maybe$Just(
				{ctor: '_Tuple2', _0: label, _1: children$}));
	});
var _elm_community$graph$Graph_Tree$leaf = function (val) {
	return A2(
		_elm_community$graph$Graph_Tree$inner,
		val,
		_elm_lang$core$Native_List.fromArray(
			[]));
};
var _elm_community$graph$Graph_Tree$unfoldTree = F2(
	function (next, seed) {
		var _p15 = next(seed);
		var label = _p15._0;
		var seeds = _p15._1;
		return A2(
			_elm_community$graph$Graph_Tree$inner,
			label,
			A2(
				_elm_lang$core$List$map,
				_elm_community$graph$Graph_Tree$unfoldTree(next),
				seeds));
	});
var _elm_community$graph$Graph_Tree$unfoldForest = F2(
	function (next, seeds) {
		return A2(
			_elm_lang$core$List$map,
			_elm_community$graph$Graph_Tree$unfoldTree(next),
			seeds);
	});

var _elm_lang$core$Native_Bitwise = function() {

return {
	and: F2(function and(a, b) { return a & b; }),
	or: F2(function or(a, b) { return a | b; }),
	xor: F2(function xor(a, b) { return a ^ b; }),
	complement: function complement(a) { return ~a; },
	shiftLeft: F2(function sll(a, offset) { return a << offset; }),
	shiftRightArithmatic: F2(function sra(a, offset) { return a >> offset; }),
	shiftRightLogical: F2(function srl(a, offset) { return a >>> offset; })
};

}();

var _elm_lang$core$Bitwise$shiftRightLogical = _elm_lang$core$Native_Bitwise.shiftRightLogical;
var _elm_lang$core$Bitwise$shiftRight = _elm_lang$core$Native_Bitwise.shiftRightArithmatic;
var _elm_lang$core$Bitwise$shiftLeft = _elm_lang$core$Native_Bitwise.shiftLeft;
var _elm_lang$core$Bitwise$complement = _elm_lang$core$Native_Bitwise.complement;
var _elm_lang$core$Bitwise$xor = _elm_lang$core$Native_Bitwise.xor;
var _elm_lang$core$Bitwise$or = _elm_lang$core$Native_Bitwise.or;
var _elm_lang$core$Bitwise$and = _elm_lang$core$Native_Bitwise.and;

var _elm_community$intdict$IntDict$combineBits = F3(
	function (a, b, mask) {
		return A2(
			_elm_lang$core$Bitwise$or,
			A2(
				_elm_lang$core$Bitwise$and,
				a,
				_elm_lang$core$Bitwise$complement(mask)),
			A2(_elm_lang$core$Bitwise$and, b, mask));
	});
var _elm_community$intdict$IntDict$foldr = F3(
	function (f, acc, dict) {
		foldr:
		while (true) {
			var _p0 = dict;
			switch (_p0.ctor) {
				case 'Empty':
					return acc;
				case 'Leaf':
					var _p1 = _p0._0;
					return A3(f, _p1.key, _p1.value, acc);
				default:
					var _p2 = _p0._0;
					var _v1 = f,
						_v2 = A3(_elm_community$intdict$IntDict$foldr, f, acc, _p2.right),
						_v3 = _p2.left;
					f = _v1;
					acc = _v2;
					dict = _v3;
					continue foldr;
			}
		}
	});
var _elm_community$intdict$IntDict$keys = function (dict) {
	return A3(
		_elm_community$intdict$IntDict$foldr,
		F3(
			function (key, value, keyList) {
				return A2(_elm_lang$core$List_ops['::'], key, keyList);
			}),
		_elm_lang$core$Native_List.fromArray(
			[]),
		dict);
};
var _elm_community$intdict$IntDict$values = function (dict) {
	return A3(
		_elm_community$intdict$IntDict$foldr,
		F3(
			function (key, value, valueList) {
				return A2(_elm_lang$core$List_ops['::'], value, valueList);
			}),
		_elm_lang$core$Native_List.fromArray(
			[]),
		dict);
};
var _elm_community$intdict$IntDict$toList = function (dict) {
	return A3(
		_elm_community$intdict$IntDict$foldr,
		F3(
			function (key, value, list) {
				return A2(
					_elm_lang$core$List_ops['::'],
					{ctor: '_Tuple2', _0: key, _1: value},
					list);
			}),
		_elm_lang$core$Native_List.fromArray(
			[]),
		dict);
};
var _elm_community$intdict$IntDict$toString$ = function (dict) {
	return A2(
		_elm_lang$core$Basics_ops['++'],
		'IntDict.fromList ',
		_elm_lang$core$Basics$toString(
			_elm_community$intdict$IntDict$toList(dict)));
};
var _elm_community$intdict$IntDict$foldl = F3(
	function (f, acc, dict) {
		foldl:
		while (true) {
			var _p3 = dict;
			switch (_p3.ctor) {
				case 'Empty':
					return acc;
				case 'Leaf':
					var _p4 = _p3._0;
					return A3(f, _p4.key, _p4.value, acc);
				default:
					var _p5 = _p3._0;
					var _v5 = f,
						_v6 = A3(_elm_community$intdict$IntDict$foldl, f, acc, _p5.left),
						_v7 = _p5.right;
					f = _v5;
					acc = _v6;
					dict = _v7;
					continue foldl;
			}
		}
	});
var _elm_community$intdict$IntDict$findMax = function (dict) {
	findMax:
	while (true) {
		var _p6 = dict;
		switch (_p6.ctor) {
			case 'Empty':
				return _elm_lang$core$Maybe$Nothing;
			case 'Leaf':
				var _p7 = _p6._0;
				return _elm_lang$core$Maybe$Just(
					{ctor: '_Tuple2', _0: _p7.key, _1: _p7.value});
			default:
				var _v9 = _p6._0.right;
				dict = _v9;
				continue findMax;
		}
	}
};
var _elm_community$intdict$IntDict$findMin = function (dict) {
	findMin:
	while (true) {
		var _p8 = dict;
		switch (_p8.ctor) {
			case 'Empty':
				return _elm_lang$core$Maybe$Nothing;
			case 'Leaf':
				var _p9 = _p8._0;
				return _elm_lang$core$Maybe$Just(
					{ctor: '_Tuple2', _0: _p9.key, _1: _p9.value});
			default:
				var _v11 = _p8._0.left;
				dict = _v11;
				continue findMin;
		}
	}
};
var _elm_community$intdict$IntDict$size = function (dict) {
	var _p10 = dict;
	switch (_p10.ctor) {
		case 'Empty':
			return 0;
		case 'Leaf':
			return 1;
		default:
			return _p10._0.size;
	}
};
var _elm_community$intdict$IntDict$isEmpty = function (dict) {
	var _p11 = dict;
	if (_p11.ctor === 'Empty') {
		return true;
	} else {
		return false;
	}
};
var _elm_community$intdict$IntDict$highestBitSet = function (n) {
	var shiftOr = F2(
		function (n$, shift) {
			return A2(
				_elm_lang$core$Bitwise$or,
				n$,
				A2(_elm_lang$core$Bitwise$shiftRightLogical, n$, shift));
		});
	var n1 = A2(shiftOr, n, 1);
	var n2 = A2(shiftOr, n1, 2);
	var n3 = A2(shiftOr, n2, 4);
	var n4 = A2(shiftOr, n3, 8);
	var n5 = A2(shiftOr, n4, 16);
	return A2(
		_elm_lang$core$Bitwise$and,
		n5,
		_elm_lang$core$Bitwise$complement(
			A2(_elm_lang$core$Bitwise$shiftRightLogical, n5, 1)));
};
var _elm_community$intdict$IntDict$signBit = _elm_community$intdict$IntDict$highestBitSet(-1);
var _elm_community$intdict$IntDict$isBranchingBitSet = F2(
	function (p, n) {
		var n$ = A2(_elm_lang$core$Bitwise$xor, n, _elm_community$intdict$IntDict$signBit);
		return !_elm_lang$core$Native_Utils.eq(
			A2(_elm_lang$core$Bitwise$and, n$, p.branchingBit),
			0);
	});
var _elm_community$intdict$IntDict$higherBitMask = function (branchingBit) {
	return _elm_lang$core$Bitwise$complement((branchingBit * 2) - 1);
};
var _elm_community$intdict$IntDict$prefixMatches = F2(
	function (p, n) {
		return _elm_lang$core$Native_Utils.eq(
			A2(
				_elm_lang$core$Bitwise$and,
				n,
				_elm_community$intdict$IntDict$higherBitMask(p.branchingBit)),
			p.prefixBits);
	});
var _elm_community$intdict$IntDict$get = F2(
	function (key, dict) {
		get:
		while (true) {
			var _p12 = dict;
			switch (_p12.ctor) {
				case 'Empty':
					return _elm_lang$core$Maybe$Nothing;
				case 'Leaf':
					var _p13 = _p12._0;
					return _elm_lang$core$Native_Utils.eq(_p13.key, key) ? _elm_lang$core$Maybe$Just(_p13.value) : _elm_lang$core$Maybe$Nothing;
				default:
					var _p14 = _p12._0;
					if (_elm_lang$core$Basics$not(
						A2(_elm_community$intdict$IntDict$prefixMatches, _p14.prefix, key))) {
						return _elm_lang$core$Maybe$Nothing;
					} else {
						if (A2(_elm_community$intdict$IntDict$isBranchingBitSet, _p14.prefix, key)) {
							var _v15 = key,
								_v16 = _p14.right;
							key = _v15;
							dict = _v16;
							continue get;
						} else {
							var _v17 = key,
								_v18 = _p14.left;
							key = _v17;
							dict = _v18;
							continue get;
						}
					}
			}
		}
	});
var _elm_community$intdict$IntDict$member = F2(
	function (key, dict) {
		var _p15 = A2(_elm_community$intdict$IntDict$get, key, dict);
		if (_p15.ctor === 'Just') {
			return true;
		} else {
			return false;
		}
	});
var _elm_community$intdict$IntDict$lcp = F2(
	function (x, y) {
		var diff = A2(_elm_lang$core$Bitwise$xor, x, y);
		var branchingBit = _elm_community$intdict$IntDict$highestBitSet(diff);
		var mask = _elm_community$intdict$IntDict$higherBitMask(branchingBit);
		var prefixBits = A2(_elm_lang$core$Bitwise$and, x, mask);
		return {prefixBits: prefixBits, branchingBit: branchingBit};
	});
var _elm_community$intdict$IntDict$isValidKey = function (k) {
	return _elm_lang$core$Native_Utils.eq(
		A2(_elm_lang$core$Bitwise$or, k, 0),
		k);
};
var _elm_community$intdict$IntDict$KeyPrefix = F2(
	function (a, b) {
		return {prefixBits: a, branchingBit: b};
	});
var _elm_community$intdict$IntDict$InnerType = F4(
	function (a, b, c, d) {
		return {prefix: a, left: b, right: c, size: d};
	});
var _elm_community$intdict$IntDict$Inner = function (a) {
	return {ctor: 'Inner', _0: a};
};
var _elm_community$intdict$IntDict$inner = F3(
	function (p, l, r) {
		var _p16 = {ctor: '_Tuple2', _0: l, _1: r};
		if (_p16._0.ctor === 'Empty') {
			return r;
		} else {
			if (_p16._1.ctor === 'Empty') {
				return l;
			} else {
				return _elm_community$intdict$IntDict$Inner(
					{
						prefix: p,
						left: l,
						right: r,
						size: _elm_community$intdict$IntDict$size(l) + _elm_community$intdict$IntDict$size(r)
					});
			}
		}
	});
var _elm_community$intdict$IntDict$Leaf = function (a) {
	return {ctor: 'Leaf', _0: a};
};
var _elm_community$intdict$IntDict$leaf = F2(
	function (k, v) {
		return _elm_community$intdict$IntDict$Leaf(
			{key: k, value: v});
	});
var _elm_community$intdict$IntDict$singleton = F2(
	function (key, value) {
		return A2(_elm_community$intdict$IntDict$leaf, key, value);
	});
var _elm_community$intdict$IntDict$Empty = {ctor: 'Empty'};
var _elm_community$intdict$IntDict$empty = _elm_community$intdict$IntDict$Empty;
var _elm_community$intdict$IntDict$update = F3(
	function (key, alter, dict) {
		var join = F2(
			function (_p18, _p17) {
				var _p19 = _p18;
				var _p23 = _p19._1;
				var _p20 = _p17;
				var _p22 = _p20._1;
				var _p21 = _p20._0;
				var prefix = A2(_elm_community$intdict$IntDict$lcp, _p19._0, _p21);
				return A2(_elm_community$intdict$IntDict$isBranchingBitSet, prefix, _p21) ? A3(_elm_community$intdict$IntDict$inner, prefix, _p23, _p22) : A3(_elm_community$intdict$IntDict$inner, prefix, _p22, _p23);
			});
		var alteredNode = function (v) {
			var _p24 = alter(v);
			if (_p24.ctor === 'Just') {
				return A2(_elm_community$intdict$IntDict$leaf, key, _p24._0);
			} else {
				return _elm_community$intdict$IntDict$empty;
			}
		};
		var _p25 = dict;
		switch (_p25.ctor) {
			case 'Empty':
				return alteredNode(_elm_lang$core$Maybe$Nothing);
			case 'Leaf':
				var _p26 = _p25._0;
				return _elm_lang$core$Native_Utils.eq(_p26.key, key) ? alteredNode(
					_elm_lang$core$Maybe$Just(_p26.value)) : A2(
					join,
					{
						ctor: '_Tuple2',
						_0: key,
						_1: alteredNode(_elm_lang$core$Maybe$Nothing)
					},
					{ctor: '_Tuple2', _0: _p26.key, _1: dict});
			default:
				var _p27 = _p25._0;
				return A2(_elm_community$intdict$IntDict$prefixMatches, _p27.prefix, key) ? (A2(_elm_community$intdict$IntDict$isBranchingBitSet, _p27.prefix, key) ? A3(
					_elm_community$intdict$IntDict$inner,
					_p27.prefix,
					_p27.left,
					A3(_elm_community$intdict$IntDict$update, key, alter, _p27.right)) : A3(
					_elm_community$intdict$IntDict$inner,
					_p27.prefix,
					A3(_elm_community$intdict$IntDict$update, key, alter, _p27.left),
					_p27.right)) : A2(
					join,
					{
						ctor: '_Tuple2',
						_0: key,
						_1: alteredNode(_elm_lang$core$Maybe$Nothing)
					},
					{ctor: '_Tuple2', _0: _p27.prefix.prefixBits, _1: dict});
		}
	});
var _elm_community$intdict$IntDict$insert = F3(
	function (key, value, dict) {
		return A3(
			_elm_community$intdict$IntDict$update,
			key,
			_elm_lang$core$Basics$always(
				_elm_lang$core$Maybe$Just(value)),
			dict);
	});
var _elm_community$intdict$IntDict$remove = F2(
	function (key, dict) {
		return A3(
			_elm_community$intdict$IntDict$update,
			key,
			_elm_lang$core$Basics$always(_elm_lang$core$Maybe$Nothing),
			dict);
	});
var _elm_community$intdict$IntDict$filter = F2(
	function (predicate, dict) {
		var add = F3(
			function (k, v, d) {
				return A2(predicate, k, v) ? A3(_elm_community$intdict$IntDict$insert, k, v, d) : d;
			});
		return A3(_elm_community$intdict$IntDict$foldl, add, _elm_community$intdict$IntDict$empty, dict);
	});
var _elm_community$intdict$IntDict$map = F2(
	function (f, dict) {
		var _p28 = dict;
		switch (_p28.ctor) {
			case 'Empty':
				return _elm_community$intdict$IntDict$empty;
			case 'Leaf':
				var _p29 = _p28._0;
				return A2(
					_elm_community$intdict$IntDict$leaf,
					_p29.key,
					A2(f, _p29.key, _p29.value));
			default:
				var _p30 = _p28._0;
				return A3(
					_elm_community$intdict$IntDict$inner,
					_p30.prefix,
					A2(_elm_community$intdict$IntDict$map, f, _p30.left),
					A2(_elm_community$intdict$IntDict$map, f, _p30.right));
		}
	});
var _elm_community$intdict$IntDict$partition = F2(
	function (predicate, dict) {
		var add = F3(
			function (key, value, _p31) {
				var _p32 = _p31;
				var _p34 = _p32._1;
				var _p33 = _p32._0;
				return A2(predicate, key, value) ? {
					ctor: '_Tuple2',
					_0: A3(_elm_community$intdict$IntDict$insert, key, value, _p33),
					_1: _p34
				} : {
					ctor: '_Tuple2',
					_0: _p33,
					_1: A3(_elm_community$intdict$IntDict$insert, key, value, _p34)
				};
			});
		return A3(
			_elm_community$intdict$IntDict$foldl,
			add,
			{ctor: '_Tuple2', _0: _elm_community$intdict$IntDict$empty, _1: _elm_community$intdict$IntDict$empty},
			dict);
	});
var _elm_community$intdict$IntDict$fromList = function (pairs) {
	var insert$ = F2(
		function (_p35, dict) {
			var _p36 = _p35;
			return A3(_elm_community$intdict$IntDict$insert, _p36._0, _p36._1, dict);
		});
	return A3(_elm_lang$core$List$foldl, insert$, _elm_community$intdict$IntDict$empty, pairs);
};
var _elm_community$intdict$IntDict$Right = {ctor: 'Right'};
var _elm_community$intdict$IntDict$Left = {ctor: 'Left'};
var _elm_community$intdict$IntDict$Disjunct = F2(
	function (a, b) {
		return {ctor: 'Disjunct', _0: a, _1: b};
	});
var _elm_community$intdict$IntDict$Parent = F2(
	function (a, b) {
		return {ctor: 'Parent', _0: a, _1: b};
	});
var _elm_community$intdict$IntDict$SamePrefix = {ctor: 'SamePrefix'};
var _elm_community$intdict$IntDict$determineBranchRelation = F2(
	function (l, r) {
		var childEdge = F2(
			function (prefix, c) {
				return A2(_elm_community$intdict$IntDict$isBranchingBitSet, prefix, c.prefix.prefixBits) ? _elm_community$intdict$IntDict$Right : _elm_community$intdict$IntDict$Left;
			});
		var rp = r.prefix;
		var lp = l.prefix;
		var mask = _elm_community$intdict$IntDict$highestBitSet(
			A2(_elm_lang$core$Basics$max, lp.branchingBit, rp.branchingBit));
		var modifiedRightPrefix = A3(
			_elm_community$intdict$IntDict$combineBits,
			rp.prefixBits,
			_elm_lang$core$Bitwise$complement(lp.prefixBits),
			mask);
		var prefix = A2(_elm_community$intdict$IntDict$lcp, lp.prefixBits, modifiedRightPrefix);
		return _elm_lang$core$Native_Utils.eq(lp, rp) ? _elm_community$intdict$IntDict$SamePrefix : (_elm_lang$core$Native_Utils.eq(prefix, lp) ? A2(
			_elm_community$intdict$IntDict$Parent,
			_elm_community$intdict$IntDict$Left,
			A2(childEdge, l.prefix, r)) : (_elm_lang$core$Native_Utils.eq(prefix, rp) ? A2(
			_elm_community$intdict$IntDict$Parent,
			_elm_community$intdict$IntDict$Right,
			A2(childEdge, r.prefix, l)) : A2(
			_elm_community$intdict$IntDict$Disjunct,
			prefix,
			A2(childEdge, prefix, l))));
	});
var _elm_community$intdict$IntDict$uniteWith = F3(
	function (merger, l, r) {
		var mergeWith = F3(
			function (key, left, right) {
				var _p37 = {ctor: '_Tuple2', _0: left, _1: right};
				if (_p37._0.ctor === 'Just') {
					if (_p37._1.ctor === 'Just') {
						return _elm_lang$core$Maybe$Just(
							A3(merger, key, _p37._0._0, _p37._1._0));
					} else {
						return left;
					}
				} else {
					if (_p37._1.ctor === 'Just') {
						return right;
					} else {
						return _elm_lang$core$Native_Utils.crashCase(
							'IntDict',
							{
								start: {line: 426, column: 7},
								end: {line: 431, column: 144}
							},
							_p37)('IntDict.uniteWith: mergeWith was called with 2 Nothings. This is a bug in the implementation, please file a bug report!');
					}
				}
			});
		var _p39 = {ctor: '_Tuple2', _0: l, _1: r};
		_v29_2:
		do {
			_v29_1:
			do {
				switch (_p39._0.ctor) {
					case 'Empty':
						return r;
					case 'Leaf':
						switch (_p39._1.ctor) {
							case 'Empty':
								break _v29_1;
							case 'Leaf':
								break _v29_2;
							default:
								break _v29_2;
						}
					default:
						switch (_p39._1.ctor) {
							case 'Empty':
								break _v29_1;
							case 'Leaf':
								var _p41 = _p39._1._0;
								return A3(
									_elm_community$intdict$IntDict$update,
									_p41.key,
									function (l$) {
										return A3(
											mergeWith,
											_p41.key,
											l$,
											_elm_lang$core$Maybe$Just(_p41.value));
									},
									l);
							default:
								var _p44 = _p39._1._0;
								var _p43 = _p39._0._0;
								var _p42 = A2(_elm_community$intdict$IntDict$determineBranchRelation, _p43, _p44);
								switch (_p42.ctor) {
									case 'SamePrefix':
										return A3(
											_elm_community$intdict$IntDict$inner,
											_p43.prefix,
											A3(_elm_community$intdict$IntDict$uniteWith, merger, _p43.left, _p44.left),
											A3(_elm_community$intdict$IntDict$uniteWith, merger, _p43.right, _p44.right));
									case 'Parent':
										if (_p42._0.ctor === 'Left') {
											if (_p42._1.ctor === 'Right') {
												return A3(
													_elm_community$intdict$IntDict$inner,
													_p43.prefix,
													_p43.left,
													A3(_elm_community$intdict$IntDict$uniteWith, merger, _p43.right, r));
											} else {
												return A3(
													_elm_community$intdict$IntDict$inner,
													_p43.prefix,
													A3(_elm_community$intdict$IntDict$uniteWith, merger, _p43.left, r),
													_p43.right);
											}
										} else {
											if (_p42._1.ctor === 'Right') {
												return A3(
													_elm_community$intdict$IntDict$inner,
													_p44.prefix,
													_p44.left,
													A3(_elm_community$intdict$IntDict$uniteWith, merger, l, _p44.right));
											} else {
												return A3(
													_elm_community$intdict$IntDict$inner,
													_p44.prefix,
													A3(_elm_community$intdict$IntDict$uniteWith, merger, l, _p44.left),
													_p44.right);
											}
										}
									default:
										if (_p42._1.ctor === 'Left') {
											return A3(_elm_community$intdict$IntDict$inner, _p42._0, l, r);
										} else {
											return A3(_elm_community$intdict$IntDict$inner, _p42._0, r, l);
										}
								}
						}
				}
			} while(false);
			return l;
		} while(false);
		var _p40 = _p39._0._0;
		return A3(
			_elm_community$intdict$IntDict$update,
			_p40.key,
			function (r$) {
				return A3(
					mergeWith,
					_p40.key,
					_elm_lang$core$Maybe$Just(_p40.value),
					r$);
			},
			r);
	});
var _elm_community$intdict$IntDict$union = _elm_community$intdict$IntDict$uniteWith(
	F3(
		function (key, old, $new) {
			return old;
		}));
var _elm_community$intdict$IntDict$intersect = F2(
	function (l, r) {
		intersect:
		while (true) {
			var _p45 = {ctor: '_Tuple2', _0: l, _1: r};
			_v31_2:
			do {
				_v31_1:
				do {
					switch (_p45._0.ctor) {
						case 'Empty':
							return _elm_community$intdict$IntDict$Empty;
						case 'Leaf':
							switch (_p45._1.ctor) {
								case 'Empty':
									break _v31_1;
								case 'Leaf':
									break _v31_2;
								default:
									break _v31_2;
							}
						default:
							switch (_p45._1.ctor) {
								case 'Empty':
									break _v31_1;
								case 'Leaf':
									var _p47 = _p45._1._0;
									var _p46 = A2(_elm_community$intdict$IntDict$get, _p47.key, l);
									if (_p46.ctor === 'Just') {
										return A2(_elm_community$intdict$IntDict$leaf, _p47.key, _p46._0);
									} else {
										return _elm_community$intdict$IntDict$Empty;
									}
								default:
									var _p50 = _p45._1._0;
									var _p49 = _p45._0._0;
									var _p48 = A2(_elm_community$intdict$IntDict$determineBranchRelation, _p49, _p50);
									switch (_p48.ctor) {
										case 'SamePrefix':
											return A3(
												_elm_community$intdict$IntDict$inner,
												_p49.prefix,
												A2(_elm_community$intdict$IntDict$intersect, _p49.left, _p50.left),
												A2(_elm_community$intdict$IntDict$intersect, _p49.right, _p50.right));
										case 'Parent':
											if (_p48._0.ctor === 'Left') {
												if (_p48._1.ctor === 'Right') {
													var _v34 = _p49.right,
														_v35 = r;
													l = _v34;
													r = _v35;
													continue intersect;
												} else {
													var _v36 = _p49.left,
														_v37 = r;
													l = _v36;
													r = _v37;
													continue intersect;
												}
											} else {
												if (_p48._1.ctor === 'Right') {
													var _v38 = l,
														_v39 = _p50.right;
													l = _v38;
													r = _v39;
													continue intersect;
												} else {
													var _v40 = l,
														_v41 = _p50.left;
													l = _v40;
													r = _v41;
													continue intersect;
												}
											}
										default:
											return _elm_community$intdict$IntDict$Empty;
									}
							}
					}
				} while(false);
				return _elm_community$intdict$IntDict$Empty;
			} while(false);
			return A2(_elm_community$intdict$IntDict$member, _p45._0._0.key, r) ? l : _elm_community$intdict$IntDict$Empty;
		}
	});
var _elm_community$intdict$IntDict$diff = F2(
	function (l, r) {
		diff:
		while (true) {
			var _p51 = {ctor: '_Tuple2', _0: l, _1: r};
			_v42_2:
			do {
				_v42_1:
				do {
					switch (_p51._0.ctor) {
						case 'Empty':
							return _elm_community$intdict$IntDict$Empty;
						case 'Leaf':
							switch (_p51._1.ctor) {
								case 'Empty':
									break _v42_1;
								case 'Leaf':
									break _v42_2;
								default:
									break _v42_2;
							}
						default:
							switch (_p51._1.ctor) {
								case 'Empty':
									break _v42_1;
								case 'Leaf':
									return A2(_elm_community$intdict$IntDict$remove, _p51._1._0.key, l);
								default:
									var _p54 = _p51._1._0;
									var _p53 = _p51._0._0;
									var _p52 = A2(_elm_community$intdict$IntDict$determineBranchRelation, _p53, _p54);
									switch (_p52.ctor) {
										case 'SamePrefix':
											return A3(
												_elm_community$intdict$IntDict$inner,
												_p53.prefix,
												A2(_elm_community$intdict$IntDict$diff, _p53.left, _p54.left),
												A2(_elm_community$intdict$IntDict$diff, _p53.right, _p54.right));
										case 'Parent':
											if (_p52._0.ctor === 'Left') {
												if (_p52._1.ctor === 'Left') {
													return A3(
														_elm_community$intdict$IntDict$inner,
														_p53.prefix,
														A2(_elm_community$intdict$IntDict$diff, _p53.left, r),
														_p53.right);
												} else {
													return A3(
														_elm_community$intdict$IntDict$inner,
														_p53.prefix,
														_p53.left,
														A2(_elm_community$intdict$IntDict$diff, _p53.right, r));
												}
											} else {
												if (_p52._1.ctor === 'Left') {
													var _v44 = l,
														_v45 = _p54.left;
													l = _v44;
													r = _v45;
													continue diff;
												} else {
													var _v46 = l,
														_v47 = _p54.right;
													l = _v46;
													r = _v47;
													continue diff;
												}
											}
										default:
											return l;
									}
							}
					}
				} while(false);
				return l;
			} while(false);
			return A2(_elm_community$intdict$IntDict$member, _p51._0._0.key, r) ? _elm_community$intdict$IntDict$Empty : l;
		}
	});
var _elm_community$intdict$IntDict$merge = F6(
	function (left, both, right, l, r, acc) {
		var m = A3(_elm_community$intdict$IntDict$merge, left, both, right);
		var _p55 = {ctor: '_Tuple2', _0: l, _1: r};
		_v48_2:
		do {
			_v48_1:
			do {
				switch (_p55._0.ctor) {
					case 'Empty':
						return A3(_elm_community$intdict$IntDict$foldl, right, acc, r);
					case 'Leaf':
						switch (_p55._1.ctor) {
							case 'Empty':
								break _v48_1;
							case 'Leaf':
								break _v48_2;
							default:
								break _v48_2;
						}
					default:
						switch (_p55._1.ctor) {
							case 'Empty':
								break _v48_1;
							case 'Leaf':
								var _p59 = _p55._1._0;
								var _p58 = A2(_elm_community$intdict$IntDict$get, _p59.key, l);
								if (_p58.ctor === 'Nothing') {
									return A3(
										m,
										l,
										_elm_community$intdict$IntDict$Empty,
										A3(right, _p59.key, _p59.value, acc));
								} else {
									return A3(
										m,
										A2(_elm_community$intdict$IntDict$remove, _p59.key, l),
										_elm_community$intdict$IntDict$Empty,
										A4(both, _p59.key, _p58._0, _p59.value, acc));
								}
							default:
								var _p62 = _p55._1._0;
								var _p61 = _p55._0._0;
								var _p60 = A2(_elm_community$intdict$IntDict$determineBranchRelation, _p61, _p62);
								switch (_p60.ctor) {
									case 'SamePrefix':
										return A3(
											m,
											_p61.right,
											_p62.right,
											A3(m, _p61.left, _p62.left, acc));
									case 'Parent':
										if (_p60._0.ctor === 'Left') {
											if (_p60._1.ctor === 'Left') {
												return A3(
													m,
													_p61.right,
													_elm_community$intdict$IntDict$Empty,
													A3(m, _p61.left, r, acc));
											} else {
												return A3(
													m,
													_p61.right,
													r,
													A3(m, _p61.left, _elm_community$intdict$IntDict$Empty, acc));
											}
										} else {
											if (_p60._1.ctor === 'Left') {
												return A3(
													m,
													_elm_community$intdict$IntDict$Empty,
													_p62.right,
													A3(m, l, _p62.left, acc));
											} else {
												return A3(
													m,
													l,
													_p62.right,
													A3(m, _elm_community$intdict$IntDict$Empty, _p62.left, acc));
											}
										}
									default:
										if (_p60._1.ctor === 'Left') {
											return A3(
												m,
												_elm_community$intdict$IntDict$Empty,
												r,
												A3(m, l, _elm_community$intdict$IntDict$Empty, acc));
										} else {
											return A3(
												m,
												l,
												_elm_community$intdict$IntDict$Empty,
												A3(m, _elm_community$intdict$IntDict$Empty, r, acc));
										}
								}
						}
				}
			} while(false);
			return A3(_elm_community$intdict$IntDict$foldl, left, acc, l);
		} while(false);
		var _p57 = _p55._0._0;
		var _p56 = A2(_elm_community$intdict$IntDict$get, _p57.key, r);
		if (_p56.ctor === 'Nothing') {
			return A3(
				m,
				_elm_community$intdict$IntDict$Empty,
				r,
				A3(left, _p57.key, _p57.value, acc));
		} else {
			return A3(
				m,
				_elm_community$intdict$IntDict$Empty,
				A2(_elm_community$intdict$IntDict$remove, _p57.key, r),
				A4(both, _p57.key, _p57.value, _p56._0, acc));
		}
	});

var _evancz$focus$Focus$update = F3(
	function (_p0, f, big) {
		var _p1 = _p0;
		return A2(_p1._0.update, f, big);
	});
var _evancz$focus$Focus$set = F3(
	function (_p2, small, big) {
		var _p3 = _p2;
		return A2(
			_p3._0.update,
			_elm_lang$core$Basics$always(small),
			big);
	});
var _evancz$focus$Focus$get = F2(
	function (_p4, big) {
		var _p5 = _p4;
		return _p5._0.get(big);
	});
var _evancz$focus$Focus$Focus = function (a) {
	return {ctor: 'Focus', _0: a};
};
var _evancz$focus$Focus$create = F2(
	function (get, update) {
		return _evancz$focus$Focus$Focus(
			{get: get, update: update});
	});
var _evancz$focus$Focus_ops = _evancz$focus$Focus_ops || {};
_evancz$focus$Focus_ops['=>'] = F2(
	function (_p7, _p6) {
		var _p8 = _p7;
		var _p11 = _p8._0;
		var _p9 = _p6;
		var _p10 = _p9._0;
		var update = F2(
			function (f, big) {
				return A2(
					_p11.update,
					_p10.update(f),
					big);
			});
		var get = function (big) {
			return _p10.get(
				_p11.get(big));
		};
		return _evancz$focus$Focus$Focus(
			{get: get, update: update});
	});

var _elm_community$graph$Graph$ignorePath = F4(
	function (visit, path, _p0, acc) {
		var _p1 = path;
		if (_p1.ctor === '[]') {
			return _elm_lang$core$Native_Utils.crashCase(
				'Graph',
				{
					start: {line: 990, column: 3},
					end: {line: 994, column: 20}
				},
				_p1)('Graph.ignorePath: No algorithm should ever pass an empty path into this BfsNodeVisitor.');
		} else {
			return A2(visit, _p1._0, acc);
		}
	});
var _elm_community$graph$Graph$onFinish = F3(
	function (visitor, ctx, acc) {
		return {
			ctor: '_Tuple2',
			_0: acc,
			_1: visitor(ctx)
		};
	});
var _elm_community$graph$Graph$onDiscovery = F3(
	function (visitor, ctx, acc) {
		return {
			ctor: '_Tuple2',
			_0: A2(visitor, ctx, acc),
			_1: _elm_lang$core$Basics$identity
		};
	});
var _elm_community$graph$Graph$alongIncomingEdges = function (ctx) {
	return _elm_community$intdict$IntDict$keys(ctx.incoming);
};
var _elm_community$graph$Graph$alongOutgoingEdges = function (ctx) {
	return _elm_community$intdict$IntDict$keys(ctx.outgoing);
};
var _elm_community$graph$Graph$lookup = function (nodeId) {
	return A2(
		_evancz$focus$Focus$create,
		_elm_community$intdict$IntDict$get(nodeId),
		_elm_community$intdict$IntDict$update(nodeId));
};
var _elm_community$graph$Graph$outgoing = A2(
	_evancz$focus$Focus$create,
	function (_) {
		return _.outgoing;
	},
	F2(
		function (update, record) {
			return _elm_lang$core$Native_Utils.update(
				record,
				{
					outgoing: update(record.outgoing)
				});
		}));
var _elm_community$graph$Graph$incoming = A2(
	_evancz$focus$Focus$create,
	function (_) {
		return _.incoming;
	},
	F2(
		function (update, record) {
			return _elm_lang$core$Native_Utils.update(
				record,
				{
					incoming: update(record.incoming)
				});
		}));
var _elm_community$graph$Graph$node = A2(
	_evancz$focus$Focus$create,
	function (_) {
		return _.node;
	},
	F2(
		function (update, record) {
			return _elm_lang$core$Native_Utils.update(
				record,
				{
					node: update(record.node)
				});
		}));
var _elm_community$graph$Graph$to = A2(
	_evancz$focus$Focus$create,
	function (_) {
		return _.to;
	},
	F2(
		function (update, record) {
			return _elm_lang$core$Native_Utils.update(
				record,
				{
					to: update(record.to)
				});
		}));
var _elm_community$graph$Graph$from = A2(
	_evancz$focus$Focus$create,
	function (_) {
		return _.from;
	},
	F2(
		function (update, record) {
			return _elm_lang$core$Native_Utils.update(
				record,
				{
					from: update(record.from)
				});
		}));
var _elm_community$graph$Graph$label = A2(
	_evancz$focus$Focus$create,
	function (_) {
		return _.label;
	},
	F2(
		function (update, record) {
			return _elm_lang$core$Native_Utils.update(
				record,
				{
					label: update(record.label)
				});
		}));
var _elm_community$graph$Graph$id = A2(
	_evancz$focus$Focus$create,
	function (_) {
		return _.id;
	},
	F2(
		function (update, record) {
			return _elm_lang$core$Native_Utils.update(
				record,
				{
					id: update(record.id)
				});
		}));
var _elm_community$graph$Graph$applyEdgeDiff = F3(
	function (nodeId, diff, graphRep) {
		var edgeUpdateToMaybe = function (edgeUpdate) {
			var _p3 = edgeUpdate;
			if (_p3.ctor === 'Insert') {
				return _elm_lang$core$Maybe$Just(_p3._0);
			} else {
				return _elm_lang$core$Maybe$Nothing;
			}
		};
		var updateAdjacency = F3(
			function (edgeFocus, updatedId, edgeUpdate) {
				var updateLbl = A2(
					_evancz$focus$Focus$set,
					edgeFocus,
					edgeUpdateToMaybe(edgeUpdate));
				return A2(
					_elm_community$intdict$IntDict$update,
					updatedId,
					_elm_lang$core$Maybe$map(updateLbl));
			});
		var foldl$ = F3(
			function (f, dict, acc) {
				return A3(_elm_community$intdict$IntDict$foldl, f, acc, dict);
			});
		return A3(
			foldl$,
			updateAdjacency(
				A2(
					_evancz$focus$Focus_ops['=>'],
					_elm_community$graph$Graph$outgoing,
					_elm_community$graph$Graph$lookup(nodeId))),
			diff.outgoing,
			A3(
				foldl$,
				updateAdjacency(
					A2(
						_evancz$focus$Focus_ops['=>'],
						_elm_community$graph$Graph$incoming,
						_elm_community$graph$Graph$lookup(nodeId))),
				diff.incoming,
				graphRep));
	});
var _elm_community$graph$Graph$emptyDiff = {incoming: _elm_community$intdict$IntDict$empty, outgoing: _elm_community$intdict$IntDict$empty};
var _elm_community$graph$Graph$unGraph = function (graph) {
	var _p4 = graph;
	return _p4._0;
};
var _elm_community$graph$Graph$edges = function (graph) {
	var foldl$ = F3(
		function (f, dict, list) {
			return A3(_elm_community$intdict$IntDict$foldl, f, list, dict);
		});
	var prependEdges = F2(
		function (node1, ctx) {
			return A2(
				foldl$,
				F2(
					function (node2, e) {
						return F2(
							function (x, y) {
								return A2(_elm_lang$core$List_ops['::'], x, y);
							})(
							{to: node2, from: node1, label: e});
					}),
				ctx.outgoing);
		});
	return A3(
		foldl$,
		prependEdges,
		_elm_community$graph$Graph$unGraph(graph),
		_elm_lang$core$Native_List.fromArray(
			[]));
};
var _elm_community$graph$Graph$Node = F2(
	function (a, b) {
		return {id: a, label: b};
	});
var _elm_community$graph$Graph$Edge = F3(
	function (a, b, c) {
		return {from: a, to: b, label: c};
	});
var _elm_community$graph$Graph$NodeContext = F3(
	function (a, b, c) {
		return {node: a, incoming: b, outgoing: c};
	});
var _elm_community$graph$Graph$EdgeDiff = F2(
	function (a, b) {
		return {incoming: a, outgoing: b};
	});
var _elm_community$graph$Graph$Graph = function (a) {
	return {ctor: 'Graph', _0: a};
};
var _elm_community$graph$Graph$empty = _elm_community$graph$Graph$Graph(_elm_community$intdict$IntDict$empty);
var _elm_community$graph$Graph$isEmpty = function (graph) {
	return _elm_lang$core$Native_Utils.eq(graph, _elm_community$graph$Graph$empty);
};
var _elm_community$graph$Graph$fromNodesAndEdges = F2(
	function (nodes, edges) {
		var addEdge = F2(
			function (edge, rep) {
				var updateIncoming = function (ctx) {
					return _elm_lang$core$Native_Utils.update(
						ctx,
						{
							incoming: A3(_elm_community$intdict$IntDict$insert, edge.from, edge.label, ctx.incoming)
						});
				};
				var updateOutgoing = function (ctx) {
					return _elm_lang$core$Native_Utils.update(
						ctx,
						{
							outgoing: A3(_elm_community$intdict$IntDict$insert, edge.to, edge.label, ctx.outgoing)
						});
				};
				return A3(
					_elm_community$intdict$IntDict$update,
					edge.to,
					_elm_lang$core$Maybe$map(updateIncoming),
					A3(
						_elm_community$intdict$IntDict$update,
						edge.from,
						_elm_lang$core$Maybe$map(updateOutgoing),
						rep));
			});
		var nodeRep = A3(
			_elm_lang$core$List$foldl,
			function (n) {
				return A2(
					_elm_community$intdict$IntDict$insert,
					n.id,
					A3(_elm_community$graph$Graph$NodeContext, n, _elm_community$intdict$IntDict$empty, _elm_community$intdict$IntDict$empty));
			},
			_elm_community$intdict$IntDict$empty,
			nodes);
		return _elm_community$graph$Graph$Graph(
			A3(_elm_lang$core$List$foldl, addEdge, nodeRep, edges));
	});
var _elm_community$graph$Graph$fromNodeLabelsAndEdgePairs = F2(
	function (labels, edges) {
		var edges$ = A2(
			_elm_lang$core$List$map,
			function (_p5) {
				var _p6 = _p5;
				return A3(
					_elm_community$graph$Graph$Edge,
					_p6._0,
					_p6._1,
					{ctor: '_Tuple0'});
			},
			edges);
		var nodes = _elm_lang$core$Basics$snd(
			A3(
				_elm_lang$core$List$foldl,
				F2(
					function (lbl, _p7) {
						var _p8 = _p7;
						var _p9 = _p8._0;
						return {
							ctor: '_Tuple2',
							_0: _p9 + 1,
							_1: A2(
								_elm_lang$core$List_ops['::'],
								A2(_elm_community$graph$Graph$Node, _p9, lbl),
								_p8._1)
						};
					}),
				{
					ctor: '_Tuple2',
					_0: 0,
					_1: _elm_lang$core$Native_List.fromArray(
						[])
				},
				labels));
		return A2(_elm_community$graph$Graph$fromNodesAndEdges, nodes, edges$);
	});
var _elm_community$graph$Graph$graphRep = A2(
	_evancz$focus$Focus$create,
	_elm_community$graph$Graph$unGraph,
	function (update) {
		return function (_p10) {
			return _elm_community$graph$Graph$Graph(
				update(
					_elm_community$graph$Graph$unGraph(_p10)));
		};
	});
var _elm_community$graph$Graph$size = function (_p11) {
	return _elm_community$intdict$IntDict$size(
		A2(_evancz$focus$Focus$get, _elm_community$graph$Graph$graphRep, _p11));
};
var _elm_community$graph$Graph$member = function (nodeId) {
	return function (_p12) {
		return A2(
			_elm_community$intdict$IntDict$member,
			nodeId,
			A2(_evancz$focus$Focus$get, _elm_community$graph$Graph$graphRep, _p12));
	};
};
var _elm_community$graph$Graph$get = function (nodeId) {
	return _evancz$focus$Focus$get(
		A2(
			_evancz$focus$Focus_ops['=>'],
			_elm_community$graph$Graph$graphRep,
			_elm_community$graph$Graph$lookup(nodeId)));
};
var _elm_community$graph$Graph$nodeIdRange = function (graph) {
	var rep = A2(_evancz$focus$Focus$get, _elm_community$graph$Graph$graphRep, graph);
	return A2(
		_elm_lang$core$Maybe$andThen,
		_elm_community$intdict$IntDict$findMin(rep),
		function (_p13) {
			var _p14 = _p13;
			return A2(
				_elm_lang$core$Maybe$andThen,
				_elm_community$intdict$IntDict$findMax(rep),
				function (_p15) {
					var _p16 = _p15;
					return _elm_lang$core$Maybe$Just(
						{ctor: '_Tuple2', _0: _p14._0, _1: _p16._0});
				});
		});
};
var _elm_community$graph$Graph$nodes = function (_p17) {
	return A2(
		_elm_lang$core$List$map,
		function (_) {
			return _.node;
		},
		_elm_community$intdict$IntDict$values(
			A2(_evancz$focus$Focus$get, _elm_community$graph$Graph$graphRep, _p17)));
};
var _elm_community$graph$Graph$toString$ = function (graph) {
	var edgeList = _elm_community$graph$Graph$edges(graph);
	var nodeList = _elm_community$graph$Graph$nodes(graph);
	return A2(
		_elm_lang$core$Basics_ops['++'],
		'Graph.fromNodesAndEdges ',
		A2(
			_elm_lang$core$Basics_ops['++'],
			_elm_lang$core$Basics$toString(nodeList),
			A2(
				_elm_lang$core$Basics_ops['++'],
				' ',
				_elm_lang$core$Basics$toString(edgeList))));
};
var _elm_community$graph$Graph$nodeIds = function (_p18) {
	return _elm_community$intdict$IntDict$keys(
		A2(_evancz$focus$Focus$get, _elm_community$graph$Graph$graphRep, _p18));
};
var _elm_community$graph$Graph$symmetricClosure = function (edgeMerger) {
	var orderedEdgeMerger = F4(
		function (from, to, outgoing, incoming) {
			return (_elm_lang$core$Native_Utils.cmp(from, to) < 1) ? A4(edgeMerger, from, to, outgoing, incoming) : A4(edgeMerger, to, from, incoming, outgoing);
		});
	var updateContext = F2(
		function (nodeId, ctx) {
			var edges = A3(
				_elm_community$intdict$IntDict$uniteWith,
				orderedEdgeMerger(nodeId),
				ctx.outgoing,
				ctx.incoming);
			return _elm_lang$core$Native_Utils.update(
				ctx,
				{outgoing: edges, incoming: edges});
		});
	return A2(
		_evancz$focus$Focus$update,
		_elm_community$graph$Graph$graphRep,
		_elm_community$intdict$IntDict$map(updateContext));
};
var _elm_community$graph$Graph$reverseEdges = function () {
	var updateContext = F2(
		function (nodeId, ctx) {
			return _elm_lang$core$Native_Utils.update(
				ctx,
				{outgoing: ctx.incoming, incoming: ctx.outgoing});
		});
	return A2(
		_evancz$focus$Focus$update,
		_elm_community$graph$Graph$graphRep,
		_elm_community$intdict$IntDict$map(updateContext));
}();
var _elm_community$graph$Graph$Remove = function (a) {
	return {ctor: 'Remove', _0: a};
};
var _elm_community$graph$Graph$Insert = function (a) {
	return {ctor: 'Insert', _0: a};
};
var _elm_community$graph$Graph$computeEdgeDiff = F2(
	function (old, $new) {
		var collectUpdates = F3(
			function (edgeUpdate, updatedId, label) {
				var replaceUpdate = function (old) {
					var _p19 = {
						ctor: '_Tuple2',
						_0: old,
						_1: edgeUpdate(label)
					};
					if (_p19._0.ctor === 'Just') {
						if (_p19._0._0.ctor === 'Remove') {
							if (_p19._1.ctor === 'Insert') {
								var _p20 = _p19._1._0;
								return _elm_lang$core$Native_Utils.eq(_p19._0._0._0, _p20) ? _elm_lang$core$Maybe$Nothing : _elm_lang$core$Maybe$Just(
									_elm_community$graph$Graph$Insert(_p20));
							} else {
								return _elm_lang$core$Native_Utils.crashCase(
									'Graph',
									{
										start: {line: 196, column: 11},
										end: {line: 206, column: 22}
									},
									_p19)('Graph.computeEdgeDiff: Collected two removals for the same edge. This is an error in the implementation of Graph and you should file a bug report!');
							}
						} else {
							return _elm_lang$core$Native_Utils.crashCase(
								'Graph',
								{
									start: {line: 196, column: 11},
									end: {line: 206, column: 22}
								},
								_p19)('Graph.computeEdgeDiff: Collected inserts before removals. This is an error in the implementation of Graph and you should file a bug report!');
						}
					} else {
						return _elm_lang$core$Maybe$Just(_p19._1);
					}
				};
				return A2(_elm_community$intdict$IntDict$update, updatedId, replaceUpdate);
			});
		var collect = F3(
			function (edgeUpdate, adj, updates) {
				return A3(
					_elm_community$intdict$IntDict$foldl,
					collectUpdates(edgeUpdate),
					updates,
					adj);
			});
		var _p23 = {ctor: '_Tuple2', _0: old, _1: $new};
		if (_p23._0.ctor === 'Nothing') {
			if (_p23._1.ctor === 'Nothing') {
				return _elm_community$graph$Graph$emptyDiff;
			} else {
				var _p25 = _p23._1._0;
				return {
					outgoing: A3(collect, _elm_community$graph$Graph$Insert, _p25.incoming, _elm_community$intdict$IntDict$empty),
					incoming: A3(collect, _elm_community$graph$Graph$Insert, _p25.outgoing, _elm_community$intdict$IntDict$empty)
				};
			}
		} else {
			if (_p23._1.ctor === 'Nothing') {
				var _p24 = _p23._0._0;
				return {
					outgoing: A3(collect, _elm_community$graph$Graph$Remove, _p24.incoming, _elm_community$intdict$IntDict$empty),
					incoming: A3(collect, _elm_community$graph$Graph$Remove, _p24.outgoing, _elm_community$intdict$IntDict$empty)
				};
			} else {
				var _p27 = _p23._0._0;
				var _p26 = _p23._1._0;
				return _elm_lang$core$Native_Utils.eq(_p27, _p26) ? _elm_community$graph$Graph$emptyDiff : {
					outgoing: A3(
						collect,
						_elm_community$graph$Graph$Insert,
						_p26.incoming,
						A3(collect, _elm_community$graph$Graph$Remove, _p27.incoming, _elm_community$intdict$IntDict$empty)),
					incoming: A3(
						collect,
						_elm_community$graph$Graph$Insert,
						_p26.outgoing,
						A3(collect, _elm_community$graph$Graph$Remove, _p27.outgoing, _elm_community$intdict$IntDict$empty))
				};
			}
		}
	});
var _elm_community$graph$Graph$update = F2(
	function (nodeId, updater) {
		var updater$ = function (rep) {
			var filterInvalidEdges = function (ctx) {
				return _elm_community$intdict$IntDict$filter(
					F2(
						function (id, _p28) {
							return _elm_lang$core$Native_Utils.eq(id, ctx.node.id) || A2(_elm_community$intdict$IntDict$member, id, rep);
						}));
			};
			var cleanUpEdges = function (ctx) {
				return A3(
					_evancz$focus$Focus$update,
					_elm_community$graph$Graph$outgoing,
					filterInvalidEdges(ctx),
					A3(
						_evancz$focus$Focus$update,
						_elm_community$graph$Graph$incoming,
						filterInvalidEdges(ctx),
						ctx));
			};
			var old = A2(_elm_community$intdict$IntDict$get, nodeId, rep);
			var $new = A2(
				_elm_lang$core$Maybe$map,
				cleanUpEdges,
				updater(old));
			var diff = A2(_elm_community$graph$Graph$computeEdgeDiff, old, $new);
			return A3(
				_elm_community$intdict$IntDict$update,
				nodeId,
				_elm_lang$core$Basics$always($new),
				A3(_elm_community$graph$Graph$applyEdgeDiff, nodeId, diff, rep));
		};
		return A2(_evancz$focus$Focus$update, _elm_community$graph$Graph$graphRep, updater$);
	});
var _elm_community$graph$Graph$insert = F2(
	function (nodeContext, graph) {
		return A3(
			_elm_community$graph$Graph$update,
			nodeContext.node.id,
			_elm_lang$core$Basics$always(
				_elm_lang$core$Maybe$Just(nodeContext)),
			graph);
	});
var _elm_community$graph$Graph$inducedSubgraph = F2(
	function (nodeIds, graph) {
		var insertContextById = F2(
			function (nodeId, acc) {
				var _p29 = A2(_elm_community$graph$Graph$get, nodeId, graph);
				if (_p29.ctor === 'Just') {
					return A2(_elm_community$graph$Graph$insert, _p29._0, acc);
				} else {
					return acc;
				}
			});
		return A3(_elm_lang$core$List$foldl, insertContextById, _elm_community$graph$Graph$empty, nodeIds);
	});
var _elm_community$graph$Graph$remove = F2(
	function (nodeId, graph) {
		return A3(
			_elm_community$graph$Graph$update,
			nodeId,
			_elm_lang$core$Basics$always(_elm_lang$core$Maybe$Nothing),
			graph);
	});
var _elm_community$graph$Graph$fold = F3(
	function (f, acc, graph) {
		var go = F2(
			function (acc, graph1) {
				go:
				while (true) {
					var maybeContext = A3(
						_elm_lang$core$Basics$flip,
						_elm_lang$core$Maybe$andThen,
						function (id) {
							return A2(_elm_community$graph$Graph$get, id, graph);
						},
						A2(
							_elm_lang$core$Maybe$map,
							_elm_lang$core$Basics$fst,
							_elm_community$graph$Graph$nodeIdRange(graph1)));
					var _p30 = maybeContext;
					if (_p30.ctor === 'Just') {
						var _p31 = _p30._0;
						var _v11 = A2(f, _p31, acc),
							_v12 = A2(_elm_community$graph$Graph$remove, _p31.node.id, graph1);
						acc = _v11;
						graph1 = _v12;
						continue go;
					} else {
						return acc;
					}
				}
			});
		return A2(go, acc, graph);
	});
var _elm_community$graph$Graph$mapContexts = function (f) {
	return A2(
		_elm_community$graph$Graph$fold,
		function (ctx) {
			return _elm_community$graph$Graph$insert(
				f(ctx));
		},
		_elm_community$graph$Graph$empty);
};
var _elm_community$graph$Graph$mapNodes = function (f) {
	return A2(
		_elm_community$graph$Graph$fold,
		function (ctx) {
			return _elm_community$graph$Graph$insert(
				_elm_lang$core$Native_Utils.update(
					ctx,
					{
						node: {
							id: ctx.node.id,
							label: f(ctx.node.label)
						}
					}));
		},
		_elm_community$graph$Graph$empty);
};
var _elm_community$graph$Graph$mapEdges = function (f) {
	return A2(
		_elm_community$graph$Graph$fold,
		function (ctx) {
			return _elm_community$graph$Graph$insert(
				_elm_lang$core$Native_Utils.update(
					ctx,
					{
						outgoing: A2(
							_elm_community$intdict$IntDict$map,
							F2(
								function (n, e) {
									return f(e);
								}),
							ctx.outgoing),
						incoming: A2(
							_elm_community$intdict$IntDict$map,
							F2(
								function (n, e) {
									return f(e);
								}),
							ctx.incoming)
					}));
		},
		_elm_community$graph$Graph$empty);
};
var _elm_community$graph$Graph$guidedDfs = F5(
	function (selectNeighbors, visitNode, seeds, acc, graph) {
		var go = F3(
			function (seeds, acc, graph) {
				go:
				while (true) {
					var _p32 = seeds;
					if (_p32.ctor === '[]') {
						return {ctor: '_Tuple2', _0: acc, _1: graph};
					} else {
						var _p38 = _p32._1;
						var _p37 = _p32._0;
						var _p33 = A2(_elm_community$graph$Graph$get, _p37, graph);
						if (_p33.ctor === 'Nothing') {
							var _v15 = _p38,
								_v16 = acc,
								_v17 = graph;
							seeds = _v15;
							acc = _v16;
							graph = _v17;
							continue go;
						} else {
							var _p36 = _p33._0;
							var _p34 = A2(visitNode, _p36, acc);
							var accAfterDiscovery = _p34._0;
							var finishNode = _p34._1;
							var _p35 = A3(
								go,
								selectNeighbors(_p36),
								accAfterDiscovery,
								A2(_elm_community$graph$Graph$remove, _p37, graph));
							var accBeforeFinish = _p35._0;
							var graph1 = _p35._1;
							var accAfterFinish = finishNode(accBeforeFinish);
							var _v18 = _p38,
								_v19 = accAfterFinish,
								_v20 = graph1;
							seeds = _v18;
							acc = _v19;
							graph = _v20;
							continue go;
						}
					}
				}
			});
		return A3(go, seeds, acc, graph);
	});
var _elm_community$graph$Graph$dfs = F3(
	function (visitNode, acc, graph) {
		return _elm_lang$core$Basics$fst(
			A5(
				_elm_community$graph$Graph$guidedDfs,
				_elm_community$graph$Graph$alongOutgoingEdges,
				visitNode,
				_elm_community$graph$Graph$nodeIds(graph),
				acc,
				graph));
	});
var _elm_community$graph$Graph$dfsForest = F2(
	function (seeds, graph) {
		var visitNode = F2(
			function (ctx, trees) {
				return {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Native_List.fromArray(
						[]),
					_1: function (children) {
						return A2(
							_elm_lang$core$List_ops['::'],
							A2(_elm_community$graph$Graph_Tree$inner, ctx, children),
							trees);
					}
				};
			});
		return _elm_lang$core$List$reverse(
			_elm_lang$core$Basics$fst(
				A5(
					_elm_community$graph$Graph$guidedDfs,
					_elm_community$graph$Graph$alongOutgoingEdges,
					visitNode,
					seeds,
					_elm_lang$core$Native_List.fromArray(
						[]),
					graph)));
	});
var _elm_community$graph$Graph$dfsTree = F2(
	function (seed, graph) {
		var _p39 = A2(
			_elm_community$graph$Graph$dfsForest,
			_elm_lang$core$Native_List.fromArray(
				[seed]),
			graph);
		if (_p39.ctor === '[]') {
			return _elm_community$graph$Graph_Tree$empty;
		} else {
			if (_p39._1.ctor === '[]') {
				return _p39._0;
			} else {
				return _elm_lang$core$Native_Utils.crashCase(
					'Graph',
					{
						start: {line: 932, column: 3},
						end: {line: 938, column: 120}
					},
					_p39)('dfsTree: There can\'t be more than one DFS tree. This invariant is violated, please report this bug.');
			}
		}
	});
var _elm_community$graph$Graph$topologicalSort = function (graph) {
	return A2(
		_elm_lang$core$List$concatMap,
		_elm_community$graph$Graph_Tree$preOrderList,
		_elm_lang$core$List$reverse(
			A2(
				_elm_community$graph$Graph$dfsForest,
				_elm_community$graph$Graph$nodeIds(graph),
				graph)));
};
var _elm_community$graph$Graph$stronglyConnectedComponents = function (graph) {
	var timestamps = A3(
		_elm_community$graph$Graph$dfs,
		_elm_community$graph$Graph$onFinish(
			function (_p41) {
				return F2(
					function (x, y) {
						return A2(_elm_lang$core$List_ops['::'], x, y);
					})(
					function (_) {
						return _.id;
					}(
						function (_) {
							return _.node;
						}(_p41)));
			}),
		_elm_lang$core$Native_List.fromArray(
			[]),
		graph);
	var forest = A2(
		_elm_community$graph$Graph$dfsForest,
		timestamps,
		_elm_community$graph$Graph$reverseEdges(graph));
	var components = A2(
		_elm_lang$core$List$map,
		function (_p42) {
			return _elm_community$graph$Graph$reverseEdges(
				A3(
					_elm_lang$core$List$foldr,
					_elm_community$graph$Graph$insert,
					_elm_community$graph$Graph$empty,
					_elm_community$graph$Graph_Tree$preOrderList(_p42)));
		},
		forest);
	return components;
};
var _elm_community$graph$Graph$guidedBfs = F5(
	function (selectNeighbors, visitNode, seeds, acc, graph) {
		var enqueueMany = F4(
			function (distance, parentPath, nodeIds, queue) {
				return A3(
					_elm_lang$core$List$foldl,
					_imeckler$queue$Queue$push,
					queue,
					A2(
						_elm_lang$core$List$map,
						function (id) {
							return {ctor: '_Tuple3', _0: id, _1: parentPath, _2: distance};
						},
						nodeIds));
			});
		var go = F3(
			function (seeds, acc, graph) {
				go:
				while (true) {
					var _p43 = _imeckler$queue$Queue$pop(seeds);
					if (_p43.ctor === 'Nothing') {
						return {ctor: '_Tuple2', _0: acc, _1: graph};
					} else {
						var _p48 = _p43._0._1;
						var _p47 = _p43._0._0._0;
						var _p46 = _p43._0._0._2;
						var _p44 = A2(_elm_community$graph$Graph$get, _p47, graph);
						if (_p44.ctor === 'Nothing') {
							var _v24 = _p48,
								_v25 = acc,
								_v26 = graph;
							seeds = _v24;
							acc = _v25;
							graph = _v26;
							continue go;
						} else {
							var _p45 = _p44._0;
							var path = A2(_elm_lang$core$List_ops['::'], _p45, _p43._0._0._1);
							var acc$ = A3(visitNode, path, _p46, acc);
							var seeds$$ = A4(
								enqueueMany,
								_p46 + 1,
								path,
								selectNeighbors(_p45),
								_p48);
							var _v27 = seeds$$,
								_v28 = acc$,
								_v29 = A2(_elm_community$graph$Graph$remove, _p47, graph);
							seeds = _v27;
							acc = _v28;
							graph = _v29;
							continue go;
						}
					}
				}
			});
		return A3(
			go,
			A4(
				enqueueMany,
				0,
				_elm_lang$core$Native_List.fromArray(
					[]),
				seeds,
				_imeckler$queue$Queue$empty),
			acc,
			graph);
	});
var _elm_community$graph$Graph$bfs = F3(
	function (visitNode, acc, graph) {
		bfs:
		while (true) {
			var _p49 = A5(
				_elm_community$graph$Graph$guidedBfs,
				_elm_community$graph$Graph$alongOutgoingEdges,
				visitNode,
				_elm_community$graph$Graph$nodeIds(graph),
				acc,
				graph);
			var acc$ = _p49._0;
			var restgraph1 = _p49._1;
			var _p50 = _elm_community$graph$Graph$nodeIdRange(graph);
			if (_p50.ctor === 'Nothing') {
				return acc;
			} else {
				var _p51 = A5(
					_elm_community$graph$Graph$guidedBfs,
					_elm_community$graph$Graph$alongOutgoingEdges,
					visitNode,
					_elm_lang$core$Native_List.fromArray(
						[_p50._0._0]),
					acc,
					graph);
				var acc$ = _p51._0;
				var restgraph1 = _p51._1;
				var _v31 = visitNode,
					_v32 = acc$,
					_v33 = restgraph1;
				visitNode = _v31;
				acc = _v32;
				graph = _v33;
				continue bfs;
			}
		}
	});
var _elm_community$graph$Graph$heightLevels = function (graph) {
	var subtract = F2(
		function (a, b) {
			return b - a;
		});
	var decrementAndNoteSources = F3(
		function (id, _p53, _p52) {
			var _p54 = _p52;
			var _p58 = _p54._0;
			var indegrees$ = A3(
				_elm_community$intdict$IntDict$update,
				id,
				_elm_lang$core$Maybe$map(
					subtract(1)),
				_p54._1);
			var _p55 = A2(_elm_community$intdict$IntDict$get, id, indegrees$);
			if ((_p55.ctor === 'Just') && (_p55._0 === 0)) {
				var _p56 = A2(_elm_community$graph$Graph$get, id, graph);
				if (_p56.ctor === 'Just') {
					return {
						ctor: '_Tuple2',
						_0: A2(_elm_lang$core$List_ops['::'], _p56._0, _p58),
						_1: indegrees$
					};
				} else {
					return _elm_lang$core$Native_Utils.crashCase(
						'Graph',
						{
							start: {line: 1115, column: 13},
							end: {line: 1117, column: 154}
						},
						_p56)('Graph.heightLevels: Could not get a node of a graph which should be there by invariants. Please file a bug report!');
				}
			} else {
				return {ctor: '_Tuple2', _0: _p58, _1: indegrees$};
			}
		});
	var decrementIndegrees = F3(
		function (source, nextLevel, indegrees) {
			return A3(
				_elm_community$intdict$IntDict$foldl,
				decrementAndNoteSources,
				{ctor: '_Tuple2', _0: nextLevel, _1: indegrees},
				source.outgoing);
		});
	var go = F4(
		function (currentLevel, nextLevel, indegrees, graph) {
			var _p59 = {ctor: '_Tuple2', _0: currentLevel, _1: nextLevel};
			if (_p59._0.ctor === '[]') {
				if (_p59._1.ctor === '[]') {
					return _elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$core$Native_List.fromArray(
							[])
						]);
				} else {
					return A2(
						_elm_lang$core$List_ops['::'],
						_elm_lang$core$Native_List.fromArray(
							[]),
						A4(
							go,
							nextLevel,
							_elm_lang$core$Native_List.fromArray(
								[]),
							indegrees,
							graph));
				}
			} else {
				var _p63 = _p59._0._0;
				var _p60 = A3(decrementIndegrees, _p63, nextLevel, indegrees);
				var nextLevel$ = _p60._0;
				var indegrees$ = _p60._1;
				var _p61 = A4(
					go,
					_p59._0._1,
					nextLevel$,
					indegrees$,
					A2(_elm_community$graph$Graph$remove, _p63.node.id, graph));
				if (_p61.ctor === '[]') {
					return _elm_lang$core$Native_Utils.crashCase(
						'Graph',
						{
							start: {line: 1134, column: 13},
							end: {line: 1138, column: 44}
						},
						_p61)('Graph.heightLevels: Reached a branch which is impossible by invariants. Please file a bug report!');
				} else {
					return A2(
						_elm_lang$core$List_ops['::'],
						A2(_elm_lang$core$List_ops['::'], _p63, _p61._0),
						_p61._1);
				}
			}
		});
	var countIndegrees = A2(
		_elm_community$graph$Graph$fold,
		function (ctx) {
			return A2(
				_elm_community$intdict$IntDict$insert,
				ctx.node.id,
				_elm_community$intdict$IntDict$size(ctx.incoming));
		},
		_elm_community$intdict$IntDict$empty);
	var sources = A3(
		_elm_community$graph$Graph$fold,
		F2(
			function (ctx, acc) {
				return _elm_community$intdict$IntDict$isEmpty(ctx.incoming) ? A2(_elm_lang$core$List_ops['::'], ctx, acc) : acc;
			}),
		_elm_lang$core$Native_List.fromArray(
			[]),
		graph);
	return A4(
		go,
		sources,
		_elm_lang$core$Native_List.fromArray(
			[]),
		countIndegrees(graph),
		graph);
};
var _elm_community$graph$Graph$nodeById = function (nodeId) {
	return A2(
		_evancz$focus$Focus$create,
		_elm_community$graph$Graph$get(nodeId),
		_elm_community$graph$Graph$update(nodeId));
};
var _elm_community$graph$Graph$anyNode = function () {
	var getMinId = function (_p64) {
		return A2(
			_elm_lang$core$Maybe$map,
			_elm_lang$core$Basics$fst,
			_elm_community$intdict$IntDict$findMin(
				A2(_evancz$focus$Focus$get, _elm_community$graph$Graph$graphRep, _p64)));
	};
	var get = function (graph) {
		return A2(
			_elm_lang$core$Maybe$andThen,
			getMinId(graph),
			function (id) {
				return A2(
					_evancz$focus$Focus$get,
					_elm_community$graph$Graph$nodeById(id),
					graph);
			});
	};
	var update = F2(
		function (upd, graph) {
			var nodeId = A2(
				_elm_lang$core$Maybe$withDefault,
				0,
				getMinId(graph));
			return A3(
				_evancz$focus$Focus$update,
				_elm_community$graph$Graph$nodeById(nodeId),
				upd,
				graph);
		});
	return A2(_evancz$focus$Focus$create, get, update);
}();

var _concourse$atc$Grid$nodes = function (grid) {
	var _p0 = grid;
	switch (_p0.ctor) {
		case 'End':
			return _elm_lang$core$Set$empty;
		case 'Parallel':
			return A3(
				_elm_lang$core$List$foldl,
				F2(
					function (g, s) {
						return A2(
							_elm_lang$core$Set$union,
							s,
							_concourse$atc$Grid$nodes(g));
					}),
				_elm_lang$core$Set$empty,
				_p0._0);
		case 'Serial':
			return A2(
				_elm_lang$core$Set$union,
				_concourse$atc$Grid$nodes(_p0._0),
				_concourse$atc$Grid$nodes(_p0._1));
		default:
			return _elm_lang$core$Set$singleton(_p0._0.node.id);
	}
};
var _concourse$atc$Grid$terminals = function (grid) {
	var _p1 = grid;
	switch (_p1.ctor) {
		case 'End':
			return _elm_lang$core$Set$empty;
		case 'Parallel':
			var _p3 = _p1._0;
			return A2(
				_elm_lang$core$List$any,
				function (_p2) {
					return _elm_lang$core$Set$isEmpty(
						_concourse$atc$Grid$terminals(_p2));
				},
				_p3) ? _elm_lang$core$Set$empty : A3(
				_elm_lang$core$List$foldl,
				F2(
					function (g, s) {
						return A2(
							_elm_lang$core$Set$union,
							s,
							_concourse$atc$Grid$terminals(g));
					}),
				_elm_lang$core$Set$empty,
				_p3);
		case 'Serial':
			var _p4 = _p1._1;
			var bNodes = _concourse$atc$Grid$nodes(_p4);
			var bTerms = _concourse$atc$Grid$terminals(_p4);
			var aTerms = _concourse$atc$Grid$terminals(_p1._0);
			var joined = A2(_elm_lang$core$Set$union, aTerms, bTerms);
			return A2(_elm_lang$core$Set$diff, joined, bNodes);
		default:
			return _elm_lang$core$Set$fromList(
				_elm_community$intdict$IntDict$keys(_p1._0.outgoing));
	}
};
var _concourse$atc$Grid$comesDirectlyFrom = F2(
	function (up, grid) {
		comesDirectlyFrom:
		while (true) {
			var _p5 = grid;
			switch (_p5.ctor) {
				case 'End':
					return false;
				case 'Parallel':
					return A2(
						_elm_lang$core$List$any,
						_concourse$atc$Grid$comesDirectlyFrom(up),
						_p5._0);
				case 'Serial':
					var _v3 = up,
						_v4 = _p5._0;
					up = _v3;
					grid = _v4;
					continue comesDirectlyFrom;
				default:
					return A2(
						_elm_lang$core$Set$member,
						_p5._0.node.id,
						_concourse$atc$Grid$terminals(up));
			}
		}
	});
var _concourse$atc$Grid$leadsTo = F2(
	function (nc, grid) {
		var _p6 = grid;
		switch (_p6.ctor) {
			case 'End':
				return false;
			case 'Parallel':
				return A2(
					_elm_lang$core$List$any,
					_concourse$atc$Grid$leadsTo(nc),
					_p6._0);
			case 'Serial':
				return A2(_concourse$atc$Grid$leadsTo, nc, _p6._0) || A2(_concourse$atc$Grid$leadsTo, nc, _p6._1);
			default:
				return A2(_elm_community$intdict$IntDict$member, nc.node.id, _p6._0.outgoing);
		}
	});
var _concourse$atc$Grid$height = F2(
	function (nh, grid) {
		var _p7 = grid;
		switch (_p7.ctor) {
			case 'End':
				return 0;
			case 'Serial':
				return A2(
					_elm_lang$core$Basics$max,
					A2(_concourse$atc$Grid$height, nh, _p7._0),
					A2(_concourse$atc$Grid$height, nh, _p7._1));
			case 'Parallel':
				return _elm_lang$core$List$sum(
					A2(
						_elm_lang$core$List$map,
						_concourse$atc$Grid$height(nh),
						_p7._0));
			default:
				return nh(_p7._0);
		}
	});
var _concourse$atc$Grid$width = function (grid) {
	var _p8 = grid;
	switch (_p8.ctor) {
		case 'End':
			return 0;
		case 'Serial':
			return _concourse$atc$Grid$width(_p8._0) + _concourse$atc$Grid$width(_p8._1);
		case 'Parallel':
			return A2(
				_elm_lang$core$Maybe$withDefault,
				0,
				_elm_lang$core$List$maximum(
					A2(_elm_lang$core$List$map, _concourse$atc$Grid$width, _p8._0)));
		default:
			return 1;
	}
};
var _concourse$atc$Grid$showMatrix = function (m) {
	var showCell = function (c) {
		var _p9 = c;
		switch (_p9.ctor) {
			case 'MatrixSpacer':
				return '  ';
			case 'MatrixFilled':
				return '--';
			default:
				var _p10 = _p9._0;
				return (_elm_lang$core$Native_Utils.cmp(_p10.node.id, 10) < 0) ? A2(
					_elm_lang$core$Basics_ops['++'],
					' ',
					_elm_lang$core$Basics$toString(_p10.node.id)) : _elm_lang$core$Basics$toString(_p10.node.id);
		}
	};
	return A2(
		_elm_lang$core$String$join,
		'\n',
		A2(
			_elm_lang$core$List$map,
			function (r) {
				return A2(
					_elm_lang$core$String$join,
					'|',
					A2(_elm_lang$core$List$map, showCell, r));
			},
			_chendrix$elm_matrix$Matrix$toList(m)));
};
var _concourse$atc$Grid$End = {ctor: 'End'};
var _concourse$atc$Grid$Parallel = function (a) {
	return {ctor: 'Parallel', _0: a};
};
var _concourse$atc$Grid$addToStart = F2(
	function (a, b) {
		var _p11 = b;
		switch (_p11.ctor) {
			case 'End':
				return a;
			case 'Parallel':
				var _p13 = _p11._0;
				var _p12 = a;
				if (_p12.ctor === 'Parallel') {
					return _concourse$atc$Grid$Parallel(
						A2(_elm_lang$core$Basics_ops['++'], _p13, _p12._0));
				} else {
					return _concourse$atc$Grid$Parallel(
						A2(
							_elm_lang$core$Basics_ops['++'],
							_p13,
							_elm_lang$core$Native_List.fromArray(
								[a])));
				}
			default:
				var _p14 = a;
				if (_p14.ctor === 'Parallel') {
					return _concourse$atc$Grid$Parallel(
						A2(_elm_lang$core$List_ops['::'], b, _p14._0));
				} else {
					return _concourse$atc$Grid$Parallel(
						_elm_lang$core$Native_List.fromArray(
							[b, a]));
				}
		}
	});
var _concourse$atc$Grid$extractExclusiveUpstreams = F2(
	function (target, grid) {
		var _p15 = grid;
		switch (_p15.ctor) {
			case 'End':
				return {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Maybe$Just(grid),
					_1: _elm_lang$core$Native_List.fromArray(
						[])
				};
			case 'Parallel':
				var recurse = A2(
					_elm_lang$core$List$map,
					_concourse$atc$Grid$extractExclusiveUpstreams(target),
					_p15._0);
				var remainders = A2(_elm_lang$core$List$map, _elm_lang$core$Basics$fst, recurse);
				var exclusives = A2(_elm_lang$core$List$concatMap, _elm_lang$core$Basics$snd, recurse);
				return A2(
					_elm_lang$core$List$all,
					F2(
						function (x, y) {
							return _elm_lang$core$Native_Utils.eq(x, y);
						})(_elm_lang$core$Maybe$Nothing),
					remainders) ? {ctor: '_Tuple2', _0: _elm_lang$core$Maybe$Nothing, _1: exclusives} : {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Maybe$Just(
						_concourse$atc$Grid$Parallel(
							A2(_elm_lang$core$List$filterMap, _elm_lang$core$Basics$identity, remainders))),
					_1: exclusives
				};
			case 'Serial':
				var terms = _concourse$atc$Grid$terminals(grid);
				return (_elm_lang$core$Native_Utils.eq(
					_elm_lang$core$Set$size(terms),
					1) && A2(_elm_lang$core$Set$member, target.node.id, terms)) ? {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Maybe$Nothing,
					_1: _elm_lang$core$Native_List.fromArray(
						[grid])
				} : {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Maybe$Just(grid),
					_1: _elm_lang$core$Native_List.fromArray(
						[])
				};
			default:
				var _p16 = _p15._0;
				return (_elm_lang$core$Native_Utils.eq(
					_elm_community$intdict$IntDict$size(_p16.outgoing),
					1) && A2(_elm_community$intdict$IntDict$member, target.node.id, _p16.outgoing)) ? {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Maybe$Nothing,
					_1: _elm_lang$core$Native_List.fromArray(
						[grid])
				} : {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Maybe$Just(grid),
					_1: _elm_lang$core$Native_List.fromArray(
						[])
				};
		}
	});
var _concourse$atc$Grid$Serial = F2(
	function (a, b) {
		return {ctor: 'Serial', _0: a, _1: b};
	});
var _concourse$atc$Grid$addBeforeDownstream = F2(
	function (up, grid) {
		var _p17 = grid;
		switch (_p17.ctor) {
			case 'End':
				return _concourse$atc$Grid$End;
			case 'Parallel':
				return A2(_concourse$atc$Grid$comesDirectlyFrom, up, grid) ? _elm_lang$core$Native_Utils.crash(
					'Grid',
					{
						start: {line: 208, column: 9},
						end: {line: 208, column: 20}
					})('too late to add in front of Parallel') : _concourse$atc$Grid$Parallel(
					A2(
						_elm_lang$core$List$map,
						_concourse$atc$Grid$addBeforeDownstream(up),
						_p17._0));
			case 'Serial':
				var _p19 = _p17._1;
				var _p18 = _p17._0;
				return A2(_concourse$atc$Grid$comesDirectlyFrom, up, _p18) ? _elm_lang$core$Native_Utils.crash(
					'Grid',
					{
						start: {line: 214, column: 9},
						end: {line: 214, column: 20}
					})('too late to add in front of Serial') : (A2(_concourse$atc$Grid$comesDirectlyFrom, up, _p19) ? A2(
					_concourse$atc$Grid$Serial,
					A2(_concourse$atc$Grid$addToStart, up, _p18),
					_p19) : A2(
					_concourse$atc$Grid$Serial,
					_p18,
					A2(_concourse$atc$Grid$addBeforeDownstream, up, _p19)));
			default:
				return A2(_concourse$atc$Grid$comesDirectlyFrom, up, grid) ? _elm_lang$core$Native_Utils.crash(
					'Grid',
					{
						start: {line: 222, column: 9},
						end: {line: 222, column: 20}
					})('too late to add in front of Cell') : grid;
		}
	});
var _concourse$atc$Grid$checkAndAddBeforeDownstream = F2(
	function (up, grid) {
		var after = A2(_concourse$atc$Grid$addBeforeDownstream, up, grid);
		return _elm_lang$core$Native_Utils.eq(after, grid) ? _elm_lang$core$Native_Utils.crash(
			'Grid',
			{
				start: {line: 196, column: 7},
				end: {line: 196, column: 18}
			})(
			A2(
				_elm_lang$core$Basics_ops['++'],
				'failed to add: ',
				_elm_lang$core$Basics$toString(up))) : after;
	});
var _concourse$atc$Grid$Cell = function (a) {
	return {ctor: 'Cell', _0: a};
};
var _concourse$atc$Grid$addAfterMixedUpstreamsAndReinsertExclusiveOnes = F2(
	function (nc, dependent) {
		var _p20 = A2(
			_concourse$atc$Grid$extractExclusiveUpstreams,
			nc,
			_concourse$atc$Grid$Parallel(dependent));
		var remainder = _p20._0;
		var exclusives = _p20._1;
		var _p21 = {ctor: '_Tuple2', _0: remainder, _1: exclusives};
		if (_p21._0.ctor === 'Nothing') {
			if (_p21._1.ctor === '[]') {
				return _elm_lang$core$Native_Utils.crashCase(
					'Grid',
					{
						start: {line: 173, column: 5},
						end: {line: 187, column: 21}
					},
					_p21)('impossible');
			} else {
				return A2(
					_concourse$atc$Grid$Serial,
					_concourse$atc$Grid$Parallel(exclusives),
					_concourse$atc$Grid$Cell(nc));
			}
		} else {
			if (_p21._1.ctor === '[]') {
				return A2(
					_concourse$atc$Grid$Serial,
					_concourse$atc$Grid$Parallel(dependent),
					_concourse$atc$Grid$Cell(nc));
			} else {
				return A3(
					_elm_lang$core$List$foldr,
					_concourse$atc$Grid$checkAndAddBeforeDownstream,
					A2(_concourse$atc$Grid$addAfterUpstreams, nc, _p21._0._0),
					exclusives);
			}
		}
	});
var _concourse$atc$Grid$addAfterUpstreams = F2(
	function (nc, grid) {
		var _p23 = grid;
		switch (_p23.ctor) {
			case 'End':
				return _concourse$atc$Grid$End;
			case 'Parallel':
				var _p24 = A2(
					_elm_lang$core$List$partition,
					_concourse$atc$Grid$leadsTo(nc),
					_p23._0);
				var dependent = _p24._0;
				var rest = _p24._1;
				var _p25 = dependent;
				if (_p25.ctor === '[]') {
					return grid;
				} else {
					if (_p25._1.ctor === '[]') {
						return _concourse$atc$Grid$Parallel(
							A2(
								_elm_lang$core$List_ops['::'],
								A2(_concourse$atc$Grid$addAfterUpstreams, nc, _p25._0),
								rest));
					} else {
						return A2(
							_concourse$atc$Grid$addToStart,
							_concourse$atc$Grid$Parallel(rest),
							A2(_concourse$atc$Grid$addAfterMixedUpstreamsAndReinsertExclusiveOnes, nc, dependent));
					}
				}
			case 'Serial':
				var _p27 = _p23._1;
				var _p26 = _p23._0;
				return A2(_concourse$atc$Grid$leadsTo, nc, _p26) ? A2(
					_concourse$atc$Grid$Serial,
					_p26,
					A2(
						_concourse$atc$Grid$addToStart,
						_concourse$atc$Grid$Cell(nc),
						_p27)) : A2(
					_concourse$atc$Grid$Serial,
					_p26,
					A2(_concourse$atc$Grid$addAfterUpstreams, nc, _p27));
			default:
				return A2(_elm_community$intdict$IntDict$member, nc.node.id, _p23._0.outgoing) ? A2(
					_concourse$atc$Grid$Serial,
					grid,
					_concourse$atc$Grid$Cell(nc)) : grid;
		}
	});
var _concourse$atc$Grid$insert = F2(
	function (nc, grid) {
		var _p28 = _elm_community$intdict$IntDict$size(nc.incoming);
		if (_p28 === 0) {
			return A2(
				_concourse$atc$Grid$addToStart,
				_concourse$atc$Grid$Cell(nc),
				grid);
		} else {
			return A2(_concourse$atc$Grid$addAfterUpstreams, nc, grid);
		}
	});
var _concourse$atc$Grid$fromGraph = function (graph) {
	return A3(
		_elm_lang$core$List$foldl,
		_concourse$atc$Grid$insert,
		_concourse$atc$Grid$End,
		_elm_lang$core$List$concat(
			_elm_community$graph$Graph$heightLevels(graph)));
};
var _concourse$atc$Grid$MatrixFilled = {ctor: 'MatrixFilled'};
var _concourse$atc$Grid$clearHeight = F4(
	function (row, col, height, matrix) {
		clearHeight:
		while (true) {
			if (_elm_lang$core$Native_Utils.eq(height, 0)) {
				return matrix;
			} else {
				var _v18 = row,
					_v19 = col,
					_v20 = height - 1,
					_v21 = A3(
					_chendrix$elm_matrix$Matrix$set,
					{ctor: '_Tuple2', _0: row + height, _1: col},
					_concourse$atc$Grid$MatrixFilled,
					matrix);
				row = _v18;
				col = _v19;
				height = _v20;
				matrix = _v21;
				continue clearHeight;
			}
		}
	});
var _concourse$atc$Grid$MatrixSpacer = {ctor: 'MatrixSpacer'};
var _concourse$atc$Grid$MatrixNode = function (a) {
	return {ctor: 'MatrixNode', _0: a};
};
var _concourse$atc$Grid$toMatrix$ = F5(
	function (nh, row, col, matrix, grid) {
		toMatrix$:
		while (true) {
			var _p29 = grid;
			switch (_p29.ctor) {
				case 'End':
					return matrix;
				case 'Serial':
					var _p30 = _p29._0;
					var _v23 = nh,
						_v24 = row,
						_v25 = col + _concourse$atc$Grid$width(_p30),
						_v26 = A5(_concourse$atc$Grid$toMatrix$, nh, row, col, matrix, _p30),
						_v27 = _p29._1;
					nh = _v23;
					row = _v24;
					col = _v25;
					matrix = _v26;
					grid = _v27;
					continue toMatrix$;
				case 'Parallel':
					return _elm_lang$core$Basics$fst(
						A3(
							_elm_lang$core$List$foldl,
							F2(
								function (g, _p31) {
									var _p32 = _p31;
									var _p33 = _p32._1;
									return {
										ctor: '_Tuple2',
										_0: A5(_concourse$atc$Grid$toMatrix$, nh, _p33, col, _p32._0, g),
										_1: _p33 + A2(_concourse$atc$Grid$height, nh, g)
									};
								}),
							{ctor: '_Tuple2', _0: matrix, _1: row},
							_p29._0));
				default:
					var _p34 = _p29._0;
					return A3(
						_chendrix$elm_matrix$Matrix$set,
						{ctor: '_Tuple2', _0: row, _1: col},
						_concourse$atc$Grid$MatrixNode(_p34),
						A4(
							_concourse$atc$Grid$clearHeight,
							row,
							col,
							nh(_p34) - 1,
							matrix));
			}
		}
	});
var _concourse$atc$Grid$toMatrix = F2(
	function (nh, grid) {
		return A5(
			_concourse$atc$Grid$toMatrix$,
			nh,
			0,
			0,
			A3(
				_chendrix$elm_matrix$Matrix$matrix,
				A2(_concourse$atc$Grid$height, nh, grid),
				_concourse$atc$Grid$width(grid),
				_elm_lang$core$Basics$always(_concourse$atc$Grid$MatrixSpacer)),
			grid);
	});

var _concourse$atc$Job$paginationParam = function (page) {
	var _p0 = page.direction;
	if (_p0.ctor === 'Since') {
		return A2(
			_elm_lang$core$Basics_ops['++'],
			'since=',
			_elm_lang$core$Basics$toString(_p0._0));
	} else {
		return A2(
			_elm_lang$core$Basics_ops['++'],
			'until=',
			_elm_lang$core$Basics$toString(_p0._0));
	}
};
var _concourse$atc$Job$viewVersion = function (version) {
	return function (_p1) {
		return _concourse$atc$DictView$view(
			A2(
				_elm_lang$core$Dict$map,
				F2(
					function (_p2, s) {
						return _elm_lang$html$Html$text(s);
					}),
				_p1));
	}(version);
};
var _concourse$atc$Job$viewBuildOutputs = F2(
	function (model, bo) {
		return A2(
			_elm_lang$html$Html$tr,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$class('mbs pas resource fl clearfix')
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(
					_elm_lang$html$Html$td,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('resource-name mrm')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html$text(bo.resource)
						])),
					A2(
					_elm_lang$html$Html$td,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('resource-version')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							_concourse$atc$Job$viewVersion(bo.version)
						]))
				]));
	});
var _concourse$atc$Job$viewBuildInputs = F2(
	function (model, bi) {
		return A2(
			_elm_lang$html$Html$tr,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$class('mbs pas resource fl clearfix')
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(
					_elm_lang$html$Html$td,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('resource-name mrm')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html$text(bi.resource)
						])),
					A2(
					_elm_lang$html$Html$td,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('resource-version')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							_concourse$atc$Job$viewVersion(bi.version)
						]))
				]));
	});
var _concourse$atc$Job$viewBuildHeader = F2(
	function (model, b) {
		return A2(
			_elm_lang$html$Html$a,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$class(
					_concourse$atc$Concourse_BuildStatus$show(b.status)),
					_elm_lang$html$Html_Attributes$href(
					_concourse$atc$Concourse_Build$url(b))
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html$text(
					A2(_elm_lang$core$Basics_ops['++'], '#', b.name))
				]));
	});
var _concourse$atc$Job$viewPaginationBar = function (model) {
	return A2(
		_elm_lang$html$Html$div,
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html_Attributes$class('pagination fr')
			]),
		_elm_lang$core$Native_List.fromArray(
			[
				function () {
				var _p3 = model.pagination.previousPage;
				if (_p3.ctor === 'Nothing') {
					return A2(
						_elm_lang$html$Html$div,
						_elm_lang$core$Native_List.fromArray(
							[
								_elm_lang$html$Html_Attributes$class('btn-page-link disabled')
							]),
						_elm_lang$core$Native_List.fromArray(
							[
								A2(
								_elm_lang$html$Html$span,
								_elm_lang$core$Native_List.fromArray(
									[
										_elm_lang$html$Html_Attributes$class('arrow')
									]),
								_elm_lang$core$Native_List.fromArray(
									[
										A2(
										_elm_lang$html$Html$i,
										_elm_lang$core$Native_List.fromArray(
											[
												_elm_lang$html$Html_Attributes$class('fa fa-arrow-left')
											]),
										_elm_lang$core$Native_List.fromArray(
											[]))
									]))
							]));
				} else {
					return A2(
						_elm_lang$html$Html$div,
						_elm_lang$core$Native_List.fromArray(
							[
								_elm_lang$html$Html_Attributes$class('btn-page-link')
							]),
						_elm_lang$core$Native_List.fromArray(
							[
								A2(
								_elm_lang$html$Html$a,
								_elm_lang$core$Native_List.fromArray(
									[
										_elm_lang$html$Html_Attributes$class('arrow'),
										_elm_lang$html$Html_Attributes$href(
										A2(
											_elm_lang$core$Basics_ops['++'],
											'/teams/',
											A2(
												_elm_lang$core$Basics_ops['++'],
												model.jobInfo.teamName,
												A2(
													_elm_lang$core$Basics_ops['++'],
													'/pipelines/',
													A2(
														_elm_lang$core$Basics_ops['++'],
														model.jobInfo.pipelineName,
														A2(
															_elm_lang$core$Basics_ops['++'],
															'/jobs/',
															A2(
																_elm_lang$core$Basics_ops['++'],
																model.jobInfo.name,
																A2(
																	_elm_lang$core$Basics_ops['++'],
																	'?',
																	_concourse$atc$Job$paginationParam(_p3._0))))))))),
										A2(_elm_lang$html$Html_Attributes$attribute, 'aria-label', 'Previous Page')
									]),
								_elm_lang$core$Native_List.fromArray(
									[
										A2(
										_elm_lang$html$Html$i,
										_elm_lang$core$Native_List.fromArray(
											[
												_elm_lang$html$Html_Attributes$class('fa fa-arrow-left')
											]),
										_elm_lang$core$Native_List.fromArray(
											[]))
									]))
							]));
				}
			}(),
				function () {
				var _p4 = model.pagination.nextPage;
				if (_p4.ctor === 'Nothing') {
					return A2(
						_elm_lang$html$Html$div,
						_elm_lang$core$Native_List.fromArray(
							[
								_elm_lang$html$Html_Attributes$class('btn-page-link disabled')
							]),
						_elm_lang$core$Native_List.fromArray(
							[
								A2(
								_elm_lang$html$Html$span,
								_elm_lang$core$Native_List.fromArray(
									[
										_elm_lang$html$Html_Attributes$class('arrow')
									]),
								_elm_lang$core$Native_List.fromArray(
									[
										A2(
										_elm_lang$html$Html$i,
										_elm_lang$core$Native_List.fromArray(
											[
												_elm_lang$html$Html_Attributes$class('fa fa-arrow-right')
											]),
										_elm_lang$core$Native_List.fromArray(
											[]))
									]))
							]));
				} else {
					return A2(
						_elm_lang$html$Html$div,
						_elm_lang$core$Native_List.fromArray(
							[
								_elm_lang$html$Html_Attributes$class('btn-page-link')
							]),
						_elm_lang$core$Native_List.fromArray(
							[
								A2(
								_elm_lang$html$Html$a,
								_elm_lang$core$Native_List.fromArray(
									[
										_elm_lang$html$Html_Attributes$class('arrow'),
										_elm_lang$html$Html_Attributes$href(
										A2(
											_elm_lang$core$Basics_ops['++'],
											'/teams/',
											A2(
												_elm_lang$core$Basics_ops['++'],
												model.jobInfo.teamName,
												A2(
													_elm_lang$core$Basics_ops['++'],
													'/pipelines/',
													A2(
														_elm_lang$core$Basics_ops['++'],
														model.jobInfo.pipelineName,
														A2(
															_elm_lang$core$Basics_ops['++'],
															'/jobs/',
															A2(
																_elm_lang$core$Basics_ops['++'],
																model.jobInfo.name,
																A2(
																	_elm_lang$core$Basics_ops['++'],
																	'?',
																	_concourse$atc$Job$paginationParam(_p4._0))))))))),
										A2(_elm_lang$html$Html_Attributes$attribute, 'aria-label', 'Next Page')
									]),
								_elm_lang$core$Native_List.fromArray(
									[
										A2(
										_elm_lang$html$Html$i,
										_elm_lang$core$Native_List.fromArray(
											[
												_elm_lang$html$Html_Attributes$class('fa fa-arrow-right')
											]),
										_elm_lang$core$Native_List.fromArray(
											[]))
									]))
							]));
				}
			}()
			]));
};
var _concourse$atc$Job$headerBuildStatusClass = function (finishedBuild) {
	var _p5 = finishedBuild;
	if (_p5.ctor === 'Nothing') {
		return '';
	} else {
		return _concourse$atc$Concourse_BuildStatus$show(_p5._0.status);
	}
};
var _concourse$atc$Job$loadSpinner = A2(
	_elm_lang$html$Html$div,
	_elm_lang$core$Native_List.fromArray(
		[
			_elm_lang$html$Html_Attributes$class('build-step')
		]),
	_elm_lang$core$Native_List.fromArray(
		[
			A2(
			_elm_lang$html$Html$div,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$class('header')
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(
					_elm_lang$html$Html$i,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('left fa fa-fw fa-spin fa-circle-o-notch')
						]),
					_elm_lang$core$Native_List.fromArray(
						[])),
					A2(
					_elm_lang$html$Html$h3,
					_elm_lang$core$Native_List.fromArray(
						[]),
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html$text('Loading...')
						]))
				]))
		]));
var _concourse$atc$Job$viewBuildResources = F2(
	function (model, buildWithResources) {
		var outputsTable = function () {
			var _p6 = buildWithResources;
			if (_p6.ctor === 'Nothing') {
				return _concourse$atc$Job$loadSpinner;
			} else {
				return A2(
					_elm_lang$html$Html$table,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('build-resources')
						]),
					A2(
						_elm_lang$core$List$map,
						_concourse$atc$Job$viewBuildOutputs(model),
						_p6._0.resources.outputs));
			}
		}();
		var inputsTable = function () {
			var _p7 = buildWithResources;
			if (_p7.ctor === 'Nothing') {
				return _concourse$atc$Job$loadSpinner;
			} else {
				return A2(
					_elm_lang$html$Html$table,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('build-resources')
						]),
					A2(
						_elm_lang$core$List$map,
						_concourse$atc$Job$viewBuildInputs(model),
						_p7._0.resources.inputs));
			}
		}();
		return _elm_lang$core$Native_List.fromArray(
			[
				A2(
				_elm_lang$html$Html$div,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('inputs mrl')
					]),
				_elm_lang$core$Native_List.fromArray(
					[
						A2(
						_elm_lang$html$Html$div,
						_elm_lang$core$Native_List.fromArray(
							[
								_elm_lang$html$Html_Attributes$class('resource-title pbs')
							]),
						_elm_lang$core$Native_List.fromArray(
							[
								A2(
								_elm_lang$html$Html$i,
								_elm_lang$core$Native_List.fromArray(
									[
										_elm_lang$html$Html_Attributes$class('fa fa-fw fa-arrow-down prs')
									]),
								_elm_lang$core$Native_List.fromArray(
									[])),
								_elm_lang$html$Html$text('inputs')
							])),
						inputsTable
					])),
				A2(
				_elm_lang$html$Html$div,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('outputs mrl')
					]),
				_elm_lang$core$Native_List.fromArray(
					[
						A2(
						_elm_lang$html$Html$div,
						_elm_lang$core$Native_List.fromArray(
							[
								_elm_lang$html$Html_Attributes$class('resource-title pbs')
							]),
						_elm_lang$core$Native_List.fromArray(
							[
								A2(
								_elm_lang$html$Html$i,
								_elm_lang$core$Native_List.fromArray(
									[
										_elm_lang$html$Html_Attributes$class('fa fa-fw fa-arrow-up prs')
									]),
								_elm_lang$core$Native_List.fromArray(
									[])),
								_elm_lang$html$Html$text('outputs')
							])),
						outputsTable
					]))
			]);
	});
var _concourse$atc$Job$viewBuildWithResources = F2(
	function (model, lubwr) {
		return A2(
			_elm_lang$html$Html$li,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$class('js-build')
				]),
			function () {
				var buildResourcesView = A2(_concourse$atc$Job$viewBuildResources, model, lubwr.buildWithResources);
				var build = function () {
					var _p8 = lubwr.buildWithResources;
					if (_p8.ctor === 'Nothing') {
						return lubwr.nextBuild;
					} else {
						return _p8._0.build;
					}
				}();
				return _elm_lang$core$Native_List.fromArray(
					[
						A2(_concourse$atc$Job$viewBuildHeader, model, build),
						A2(
						_elm_lang$html$Html$div,
						_elm_lang$core$Native_List.fromArray(
							[
								_elm_lang$html$Html_Attributes$class('pam clearfix')
							]),
						A2(
							_elm_lang$core$List_ops['::'],
							A2(_concourse$atc$BuildDuration$view, build.duration, model.now),
							buildResourcesView))
					]);
			}());
	});
var _concourse$atc$Job$getPausedState = F2(
	function (job, pausedChanging) {
		return pausedChanging ? 'loading' : (job.paused ? 'enabled' : 'disabled');
	});
var _concourse$atc$Job$getPlayPauseLoadIcon = F2(
	function (job, pausedChanging) {
		return pausedChanging ? 'fa-circle-o-notch fa-spin' : (job.paused ? '' : 'fa-pause');
	});
var _concourse$atc$Job$isRunning = function (build) {
	return _concourse$atc$Concourse_BuildStatus$isRunning(build.status);
};
var _concourse$atc$Job$initLiveUpdatingBuildWithResources = function (nextBuild) {
	return {buildWithResources: _elm_lang$core$Maybe$Nothing, nextBuild: nextBuild};
};
var _concourse$atc$Job$addNextBuild = F2(
	function (nextBuild, buildWithResources) {
		return _elm_lang$core$Native_Utils.update(
			buildWithResources,
			{nextBuild: nextBuild});
	});
var _concourse$atc$Job$addNextBuildFromArray = F3(
	function (newBuilds, i, lubwr) {
		var _p9 = A2(_elm_lang$core$Array$get, i, newBuilds);
		if (_p9.ctor === 'Nothing') {
			return lubwr;
		} else {
			return A2(_concourse$atc$Job$addNextBuild, _p9._0, lubwr);
		}
	});
var _concourse$atc$Job$addFetchedResources = F2(
	function (resources, lubwr) {
		return _elm_lang$core$Native_Utils.update(
			lubwr,
			{
				buildWithResources: _elm_lang$core$Maybe$Just(
					{build: lubwr.nextBuild, resources: resources})
			});
	});
var _concourse$atc$Job$jobBuildsPerPage = 100;
var _concourse$atc$Job$permalink = function (builds) {
	var _p10 = _elm_lang$core$List$head(builds);
	if (_p10.ctor === 'Nothing') {
		return {
			direction: _concourse$atc$Concourse_Pagination$Since(0),
			limit: _concourse$atc$Job$jobBuildsPerPage
		};
	} else {
		return {
			direction: _concourse$atc$Concourse_Pagination$Since(_p10._0.id + 1),
			limit: _elm_lang$core$List$length(builds)
		};
	}
};
var _concourse$atc$Job$Model = F7(
	function (a, b, c, d, e, f, g) {
		return {jobInfo: a, job: b, pausedChanging: c, buildsWithResources: d, now: e, page: f, pagination: g};
	});
var _concourse$atc$Job$FetchedBuildResources = F2(
	function (a, b) {
		return {index: a, result: b};
	});
var _concourse$atc$Job$BuildWithResources = F2(
	function (a, b) {
		return {build: a, resources: b};
	});
var _concourse$atc$Job$LiveUpdatingBuildWithResources = F2(
	function (a, b) {
		return {buildWithResources: a, nextBuild: b};
	});
var _concourse$atc$Job$Flags = F5(
	function (a, b, c, d, e) {
		return {jobName: a, teamName: b, pipelineName: c, pageSince: d, pageUntil: e};
	});
var _concourse$atc$Job$PausedToggled = function (a) {
	return {ctor: 'PausedToggled', _0: a};
};
var _concourse$atc$Job$pauseJob = function (jobInfo) {
	return function (_p11) {
		return A2(
			_elm_lang$core$Platform_Cmd$map,
			_concourse$atc$Job$PausedToggled,
			A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _p11));
	}(
		_concourse$atc$Concourse_Job$pause(jobInfo));
};
var _concourse$atc$Job$unpauseJob = function (jobInfo) {
	return function (_p12) {
		return A2(
			_elm_lang$core$Platform_Cmd$map,
			_concourse$atc$Job$PausedToggled,
			A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _p12));
	}(
		_concourse$atc$Concourse_Job$unpause(jobInfo));
};
var _concourse$atc$Job$TogglePaused = {ctor: 'TogglePaused'};
var _concourse$atc$Job$view = function (model) {
	return A2(
		_elm_lang$html$Html$div,
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html_Attributes$class('with-fixed-header')
			]),
		_elm_lang$core$Native_List.fromArray(
			[
				function () {
				var _p13 = model.job;
				if (_p13.ctor === 'Nothing') {
					return _concourse$atc$Job$loadSpinner;
				} else {
					var _p14 = _p13._0;
					return A2(
						_elm_lang$html$Html$div,
						_elm_lang$core$Native_List.fromArray(
							[
								_elm_lang$html$Html_Attributes$class('fixed-header')
							]),
						_elm_lang$core$Native_List.fromArray(
							[
								A2(
								_elm_lang$html$Html$div,
								_elm_lang$core$Native_List.fromArray(
									[
										_elm_lang$html$Html_Attributes$class(
										A2(
											_elm_lang$core$Basics_ops['++'],
											'build-header ',
											_concourse$atc$Job$headerBuildStatusClass(_p14.finishedBuild)))
									]),
								_elm_lang$core$Native_List.fromArray(
									[
										A2(
										_elm_lang$html$Html$button,
										A2(
											_elm_lang$core$List$append,
											_elm_lang$core$Native_List.fromArray(
												[
													_elm_lang$html$Html_Attributes$id('job-state'),
													A2(_elm_lang$html$Html_Attributes$attribute, 'aria-label', 'Toggle Job Paused State'),
													_elm_lang$html$Html_Attributes$class(
													A2(
														_elm_lang$core$Basics_ops['++'],
														'btn-pause btn-large fl ',
														A2(_concourse$atc$Job$getPausedState, _p14, model.pausedChanging)))
												]),
											_elm_lang$core$Basics$not(model.pausedChanging) ? _elm_lang$core$Native_List.fromArray(
												[
													_elm_lang$html$Html_Events$onClick(_concourse$atc$Job$TogglePaused)
												]) : _elm_lang$core$Native_List.fromArray(
												[])),
										_elm_lang$core$Native_List.fromArray(
											[
												A2(
												_elm_lang$html$Html$i,
												_elm_lang$core$Native_List.fromArray(
													[
														_elm_lang$html$Html_Attributes$class(
														A2(
															_elm_lang$core$Basics_ops['++'],
															'fa fa-fw fa-play ',
															A2(_concourse$atc$Job$getPlayPauseLoadIcon, _p14, model.pausedChanging)))
													]),
												_elm_lang$core$Native_List.fromArray(
													[]))
											])),
										A2(
										_elm_lang$html$Html$form,
										_elm_lang$core$Native_List.fromArray(
											[
												_elm_lang$html$Html_Attributes$class('trigger-build'),
												_elm_lang$html$Html_Attributes$method('post'),
												_elm_lang$html$Html_Attributes$action(
												A2(
													_elm_lang$core$Basics_ops['++'],
													'/teams/',
													A2(
														_elm_lang$core$Basics_ops['++'],
														model.jobInfo.teamName,
														A2(
															_elm_lang$core$Basics_ops['++'],
															'/pipelines/',
															A2(
																_elm_lang$core$Basics_ops['++'],
																model.jobInfo.pipelineName,
																A2(
																	_elm_lang$core$Basics_ops['++'],
																	'/jobs/',
																	A2(_elm_lang$core$Basics_ops['++'], model.jobInfo.name, '/builds')))))))
											]),
										_elm_lang$core$Native_List.fromArray(
											[
												A2(
												_elm_lang$html$Html$button,
												_elm_lang$core$Native_List.fromArray(
													[
														_elm_lang$html$Html_Attributes$class('build-action fr'),
														_elm_lang$html$Html_Attributes$disabled(_p14.disableManualTrigger),
														A2(_elm_lang$html$Html_Attributes$attribute, 'aria-label', 'Trigger Build')
													]),
												_elm_lang$core$Native_List.fromArray(
													[
														A2(
														_elm_lang$html$Html$i,
														_elm_lang$core$Native_List.fromArray(
															[
																_elm_lang$html$Html_Attributes$class('fa fa-plus-circle')
															]),
														_elm_lang$core$Native_List.fromArray(
															[]))
													]))
											])),
										A2(
										_elm_lang$html$Html$h1,
										_elm_lang$core$Native_List.fromArray(
											[]),
										_elm_lang$core$Native_List.fromArray(
											[
												_elm_lang$html$Html$text(model.jobInfo.name)
											]))
									])),
								A2(
								_elm_lang$html$Html$div,
								_elm_lang$core$Native_List.fromArray(
									[
										_elm_lang$html$Html_Attributes$class('pagination-header')
									]),
								_elm_lang$core$Native_List.fromArray(
									[
										_concourse$atc$Job$viewPaginationBar(model),
										A2(
										_elm_lang$html$Html$h1,
										_elm_lang$core$Native_List.fromArray(
											[]),
										_elm_lang$core$Native_List.fromArray(
											[
												_elm_lang$html$Html$text('builds')
											]))
									]))
							]));
				}
			}(),
				function () {
				var _p15 = model.buildsWithResources;
				if (_p15.ctor === 'Nothing') {
					return _concourse$atc$Job$loadSpinner;
				} else {
					return A2(
						_elm_lang$html$Html$div,
						_elm_lang$core$Native_List.fromArray(
							[
								_elm_lang$html$Html_Attributes$class('scrollable-body')
							]),
						_elm_lang$core$Native_List.fromArray(
							[
								A2(
								_elm_lang$html$Html$ul,
								_elm_lang$core$Native_List.fromArray(
									[
										_elm_lang$html$Html_Attributes$class('jobs-builds-list builds-list')
									]),
								A2(
									_elm_lang$core$List$map,
									_concourse$atc$Job$viewBuildWithResources(model),
									_elm_lang$core$Array$toList(_p15._0)))
							]));
				}
			}()
			]));
};
var _concourse$atc$Job$ClockTick = function (a) {
	return {ctor: 'ClockTick', _0: a};
};
var _concourse$atc$Job$BuildResourcesFetched = function (a) {
	return {ctor: 'BuildResourcesFetched', _0: a};
};
var _concourse$atc$Job$initBuildResourcesFetched = F2(
	function (index, result) {
		return _concourse$atc$Job$BuildResourcesFetched(
			{index: index, result: result});
	});
var _concourse$atc$Job$fetchBuildResources = F2(
	function (index, buildId) {
		return function (_p16) {
			return A2(
				_elm_lang$core$Platform_Cmd$map,
				_concourse$atc$Job$initBuildResourcesFetched(index),
				A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _p16));
		}(
			_concourse$atc$Concourse_BuildResources$fetch(buildId));
	});
var _concourse$atc$Job$JobFetched = function (a) {
	return {ctor: 'JobFetched', _0: a};
};
var _concourse$atc$Job$fetchJob = F2(
	function (delay, jobInfo) {
		return function (_p17) {
			return A2(
				_elm_lang$core$Platform_Cmd$map,
				_concourse$atc$Job$JobFetched,
				A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _p17));
		}(
			A2(
				_elm_lang$core$Task$andThen,
				_elm_lang$core$Process$sleep(delay),
				_elm_lang$core$Basics$always(
					_concourse$atc$Concourse_Job$fetchJob(jobInfo))));
	});
var _concourse$atc$Job$JobBuildsFetched = function (a) {
	return {ctor: 'JobBuildsFetched', _0: a};
};
var _concourse$atc$Job$fetchJobBuilds = F3(
	function (delay, jobInfo, page) {
		return function (_p18) {
			return A2(
				_elm_lang$core$Platform_Cmd$map,
				_concourse$atc$Job$JobBuildsFetched,
				A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _p18));
		}(
			A2(
				_elm_lang$core$Task$andThen,
				_elm_lang$core$Process$sleep(delay),
				_elm_lang$core$Basics$always(
					A2(_concourse$atc$Concourse_Build$fetchJobBuilds, jobInfo, page))));
	});
var _concourse$atc$Job$init = function (flags) {
	var model = {
		jobInfo: {name: flags.jobName, teamName: flags.teamName, pipelineName: flags.pipelineName},
		job: _elm_lang$core$Maybe$Nothing,
		pausedChanging: false,
		buildsWithResources: _elm_lang$core$Maybe$Nothing,
		now: 0,
		page: {
			direction: (_elm_lang$core$Native_Utils.cmp(flags.pageUntil, 0) > 0) ? _concourse$atc$Concourse_Pagination$Until(flags.pageUntil) : _concourse$atc$Concourse_Pagination$Since(flags.pageSince),
			limit: _concourse$atc$Job$jobBuildsPerPage
		},
		pagination: {previousPage: _elm_lang$core$Maybe$Nothing, nextPage: _elm_lang$core$Maybe$Nothing}
	};
	return {
		ctor: '_Tuple2',
		_0: model,
		_1: _elm_lang$core$Platform_Cmd$batch(
			_elm_lang$core$Native_List.fromArray(
				[
					A3(
					_concourse$atc$Job$fetchJobBuilds,
					0,
					model.jobInfo,
					_elm_lang$core$Maybe$Just(model.page)),
					A2(_concourse$atc$Job$fetchJob, 0, model.jobInfo)
				]))
	};
};
var _concourse$atc$Job$handleJobBuildsFetched = F2(
	function (paginatedBuilds, model) {
		var newPage = _concourse$atc$Job$permalink(paginatedBuilds.content);
		var fetchedBuilds = _elm_lang$core$Array$fromList(paginatedBuilds.content);
		return {
			ctor: '_Tuple2',
			_0: _elm_lang$core$Native_Utils.update(
				model,
				{
					buildsWithResources: _elm_lang$core$Maybe$Just(
						function () {
							var _p19 = model.buildsWithResources;
							if (_p19.ctor === 'Nothing') {
								return A2(_elm_lang$core$Array$map, _concourse$atc$Job$initLiveUpdatingBuildWithResources, fetchedBuilds);
							} else {
								return A2(
									_elm_lang$core$Array$indexedMap,
									_concourse$atc$Job$addNextBuildFromArray(fetchedBuilds),
									_p19._0);
							}
						}()),
					page: newPage,
					pagination: paginatedBuilds.pagination
				}),
			_1: _elm_lang$core$Platform_Cmd$batch(
				A2(
					_elm_lang$core$List_ops['::'],
					A3(
						_concourse$atc$Job$fetchJobBuilds,
						5 * _elm_lang$core$Time$second,
						model.jobInfo,
						_elm_lang$core$Maybe$Just(newPage)),
					_elm_lang$core$Array$toList(
						A2(
							_elm_lang$core$Array$indexedMap,
							_concourse$atc$Job$fetchBuildResources,
							A2(
								_elm_lang$core$Array$map,
								function (_) {
									return _.id;
								},
								function () {
									var _p20 = model.buildsWithResources;
									if (_p20.ctor === 'Nothing') {
										return fetchedBuilds;
									} else {
										return A2(
											_elm_lang$core$Array$filter,
											_concourse$atc$Job$isRunning,
											A2(
												_elm_lang$core$Array$map,
												function (_) {
													return _.nextBuild;
												},
												_p20._0));
									}
								}())))))
		};
	});
var _concourse$atc$Job$Noop = {ctor: 'Noop'};
var _concourse$atc$Job$redirectToLogin = function (model) {
	return function (_p21) {
		return A2(
			_elm_lang$core$Platform_Cmd$map,
			_elm_lang$core$Basics$always(_concourse$atc$Job$Noop),
			A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _p21));
	}(
		_concourse$atc$Redirect$to('/login'));
};
var _concourse$atc$Job$update = F2(
	function (action, model) {
		var _p22 = action;
		switch (_p22.ctor) {
			case 'Noop':
				return {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none};
			case 'JobBuildsFetched':
				if (_p22._0.ctor === 'Ok') {
					return A2(_concourse$atc$Job$handleJobBuildsFetched, _p22._0._0, model);
				} else {
					return A2(
						_elm_lang$core$Debug$log,
						A2(
							_elm_lang$core$Basics_ops['++'],
							'failed to fetch builds: ',
							_elm_lang$core$Basics$toString(_p22._0._0)),
						{ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none});
				}
			case 'JobFetched':
				if (_p22._0.ctor === 'Ok') {
					return {
						ctor: '_Tuple2',
						_0: _elm_lang$core$Native_Utils.update(
							model,
							{
								job: _elm_lang$core$Maybe$Just(_p22._0._0)
							}),
						_1: A2(_concourse$atc$Job$fetchJob, 5 * _elm_lang$core$Time$second, model.jobInfo)
					};
				} else {
					return A2(
						_elm_lang$core$Debug$log,
						A2(
							_elm_lang$core$Basics_ops['++'],
							'failed to fetch job info: ',
							_elm_lang$core$Basics$toString(_p22._0._0)),
						{ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none});
				}
			case 'BuildResourcesFetched':
				var _p27 = _p22._0;
				var _p23 = _p27.result;
				if (_p23.ctor === 'Ok') {
					var _p24 = model.buildsWithResources;
					if (_p24.ctor === 'Nothing') {
						return {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none};
					} else {
						var _p26 = _p24._0;
						var _p25 = A2(_elm_lang$core$Array$get, _p27.index, _p26);
						if (_p25.ctor === 'Nothing') {
							return {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none};
						} else {
							return {
								ctor: '_Tuple2',
								_0: _elm_lang$core$Native_Utils.update(
									model,
									{
										buildsWithResources: _elm_lang$core$Maybe$Just(
											A3(
												_elm_lang$core$Array$set,
												_p27.index,
												A2(_concourse$atc$Job$addFetchedResources, _p23._0, _p25._0),
												_p26))
									}),
								_1: _elm_lang$core$Platform_Cmd$none
							};
						}
					}
				} else {
					return {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none};
				}
			case 'ClockTick':
				return {
					ctor: '_Tuple2',
					_0: _elm_lang$core$Native_Utils.update(
						model,
						{now: _p22._0}),
					_1: _elm_lang$core$Platform_Cmd$none
				};
			case 'TogglePaused':
				var _p28 = model.job;
				if (_p28.ctor === 'Nothing') {
					return {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none};
				} else {
					var _p29 = _p28._0;
					return {
						ctor: '_Tuple2',
						_0: _elm_lang$core$Native_Utils.update(
							model,
							{
								pausedChanging: true,
								job: _elm_lang$core$Maybe$Just(
									_elm_lang$core$Native_Utils.update(
										_p29,
										{
											paused: _elm_lang$core$Basics$not(_p29.paused)
										}))
							}),
						_1: _p29.paused ? _concourse$atc$Job$unpauseJob(model.jobInfo) : _concourse$atc$Job$pauseJob(model.jobInfo)
					};
				}
			default:
				if (_p22._0.ctor === 'Ok') {
					return {
						ctor: '_Tuple2',
						_0: _elm_lang$core$Native_Utils.update(
							model,
							{pausedChanging: false}),
						_1: _elm_lang$core$Platform_Cmd$none
					};
				} else {
					if ((_p22._0._0.ctor === 'BadResponse') && (_p22._0._0._0 === 401)) {
						return {
							ctor: '_Tuple2',
							_0: model,
							_1: _concourse$atc$Job$redirectToLogin(model)
						};
					} else {
						return A2(
							_elm_lang$core$Debug$log,
							A2(
								_elm_lang$core$Basics_ops['++'],
								'failed to pause/unpause job: ',
								_elm_lang$core$Basics$toString(_p22._0._0)),
							{ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none});
					}
				}
		}
	});

var _concourse$atc$JobPage$main = {
	main: _elm_lang$html$Html_App$programWithFlags(
		{
			init: _concourse$atc$Job$init,
			update: _concourse$atc$Job$update,
			view: _concourse$atc$Job$view,
			subscriptions: _elm_lang$core$Basics$always(
				A2(_elm_lang$core$Time$every, _elm_lang$core$Time$second, _concourse$atc$Job$ClockTick))
		}),
	flags: A2(
		_elm_lang$core$Json_Decode$andThen,
		A2(_elm_lang$core$Json_Decode_ops[':='], 'jobName', _elm_lang$core$Json_Decode$string),
		function (jobName) {
			return A2(
				_elm_lang$core$Json_Decode$andThen,
				A2(_elm_lang$core$Json_Decode_ops[':='], 'pageSince', _elm_lang$core$Json_Decode$int),
				function (pageSince) {
					return A2(
						_elm_lang$core$Json_Decode$andThen,
						A2(_elm_lang$core$Json_Decode_ops[':='], 'pageUntil', _elm_lang$core$Json_Decode$int),
						function (pageUntil) {
							return A2(
								_elm_lang$core$Json_Decode$andThen,
								A2(_elm_lang$core$Json_Decode_ops[':='], 'pipelineName', _elm_lang$core$Json_Decode$string),
								function (pipelineName) {
									return A2(
										_elm_lang$core$Json_Decode$andThen,
										A2(_elm_lang$core$Json_Decode_ops[':='], 'teamName', _elm_lang$core$Json_Decode$string),
										function (teamName) {
											return _elm_lang$core$Json_Decode$succeed(
												{jobName: jobName, pageSince: pageSince, pageUntil: pageUntil, pipelineName: pipelineName, teamName: teamName});
										});
								});
						});
				});
		})
};

var _concourse$atc$Login$viewOAuthButton = function (method) {
	var _p0 = method;
	if (_p0.ctor === 'BasicMethod') {
		return _elm_lang$core$Maybe$Nothing;
	} else {
		var _p1 = _p0._0;
		return _elm_lang$core$Maybe$Just(
			A2(
				_elm_lang$html$Html$a,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$href(_p1.authURL)
					]),
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html$text(
						A2(_elm_lang$core$Basics_ops['++'], 'login with ', _p1.displayName))
					])));
	}
};
var _concourse$atc$Login$viewOAuthButtons = function (methods) {
	var _p2 = A2(_elm_lang$core$List$filterMap, _concourse$atc$Login$viewOAuthButton, methods);
	if (_p2.ctor === '[]') {
		return _elm_lang$core$Maybe$Nothing;
	} else {
		return _elm_lang$core$Maybe$Just(
			A2(
				_elm_lang$html$Html$div,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('centered-contents padded-top')
					]),
				_p2));
	}
};
var _concourse$atc$Login$viewBasicAuthForm = function (methods) {
	return A2(_elm_lang$core$List$member, _concourse$atc$Concourse_AuthMethod$BasicMethod, methods) ? _elm_lang$core$Maybe$Just(
		A2(
			_elm_lang$html$Html$form,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$class('padded-top'),
					_elm_lang$html$Html_Attributes$method('post')
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					A2(
					_elm_lang$html$Html$label,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$for('basic-auth-username-input')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html$text('username')
						])),
					A2(
					_elm_lang$html$Html$div,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('input-holder')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$html$Html$input,
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html_Attributes$id('basic-auth-username-input'),
									_elm_lang$html$Html_Attributes$name('username'),
									_elm_lang$html$Html_Attributes$type$('text')
								]),
							_elm_lang$core$Native_List.fromArray(
								[]))
						])),
					A2(
					_elm_lang$html$Html$label,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$for('basic-auth-password-input')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html$text('password')
						])),
					A2(
					_elm_lang$html$Html$div,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('input-holder')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$html$Html$input,
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html_Attributes$id('basic-auth-password-input'),
									_elm_lang$html$Html_Attributes$name('password'),
									_elm_lang$html$Html_Attributes$type$('password')
								]),
							_elm_lang$core$Native_List.fromArray(
								[]))
						])),
					A2(
					_elm_lang$html$Html$div,
					_elm_lang$core$Native_List.fromArray(
						[
							_elm_lang$html$Html_Attributes$class('centered-contents')
						]),
					_elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$html$Html$button,
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html_Attributes$type$('submit')
								]),
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html$text('login')
								]))
						]))
				]))) : _elm_lang$core$Maybe$Nothing;
};
var _concourse$atc$Login$viewNoAuthButton = A2(
	_elm_lang$html$Html$form,
	_elm_lang$core$Native_List.fromArray(
		[
			_elm_lang$html$Html_Attributes$class('padded-top centered-contents'),
			_elm_lang$html$Html_Attributes$method('post')
		]),
	_elm_lang$core$Native_List.fromArray(
		[
			A2(
			_elm_lang$html$Html$button,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$type$('submit')
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html$text('login')
				]))
		]));
var _concourse$atc$Login$viewOrBar = A2(
	_elm_lang$html$Html$div,
	_elm_lang$core$Native_List.fromArray(
		[
			_elm_lang$html$Html_Attributes$class('or-bar')
		]),
	_elm_lang$core$Native_List.fromArray(
		[
			A2(
			_elm_lang$html$Html$div,
			_elm_lang$core$Native_List.fromArray(
				[]),
			_elm_lang$core$Native_List.fromArray(
				[])),
			A2(
			_elm_lang$html$Html$span,
			_elm_lang$core$Native_List.fromArray(
				[]),
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html$text('or')
				]))
		]));
var _concourse$atc$Login$teamNameStartsWithSensitive = F2(
	function (substring, team) {
		return A2(_elm_lang$core$String$startsWith, substring, team.name);
	});
var _concourse$atc$Login$teamNameStartsWith = F2(
	function (substring, team) {
		return A2(
			_elm_lang$core$String$startsWith,
			substring,
			_elm_lang$core$String$toLower(team.name));
	});
var _concourse$atc$Login$teamNameContains = F2(
	function (substring, team) {
		return A2(
			_elm_lang$core$String$contains,
			substring,
			_elm_lang$core$String$toLower(team.name));
	});
var _concourse$atc$Login$filterTeams = F2(
	function (teamFilter, teams) {
		var filteredList = A2(
			_elm_lang$core$List$filter,
			_concourse$atc$Login$teamNameContains(
				_elm_lang$core$String$toLower(teamFilter)),
			teams);
		var _p3 = A2(
			_elm_lang$core$List$partition,
			_concourse$atc$Login$teamNameStartsWith(
				_elm_lang$core$String$toLower(teamFilter)),
			filteredList);
		var startingTeams = _p3._0;
		var notStartingTeams = _p3._1;
		var _p4 = A2(
			_elm_lang$core$List$partition,
			_concourse$atc$Login$teamNameStartsWithSensitive(teamFilter),
			startingTeams);
		var caseSensitive = _p4._0;
		var notCaseSensitive = _p4._1;
		return A2(
			_elm_lang$core$Basics_ops['++'],
			caseSensitive,
			A2(_elm_lang$core$Basics_ops['++'], notCaseSensitive, notStartingTeams));
	});
var _concourse$atc$Login$assertLeftButton = F2(
	function (message, button) {
		return _elm_lang$core$Native_Utils.eq(button, 0) ? _elm_lang$core$Result$Ok(message) : _elm_lang$core$Result$Err('placeholder error, nothing is wrong');
	});
var _concourse$atc$Login$onClickPreventDefault = function (message) {
	return A3(
		_elm_lang$html$Html_Events$onWithOptions,
		'click',
		{stopPropagation: false, preventDefault: true},
		A2(
			_elm_lang$core$Json_Decode$customDecoder,
			A2(_elm_lang$core$Json_Decode_ops[':='], 'button', _elm_lang$core$Json_Decode$int),
			_concourse$atc$Login$assertLeftButton(message)));
};
var _concourse$atc$Login$viewLoading = A2(
	_elm_lang$html$Html$div,
	_elm_lang$core$Native_List.fromArray(
		[
			_elm_lang$html$Html_Attributes$class('loading')
		]),
	_elm_lang$core$Native_List.fromArray(
		[
			A2(
			_elm_lang$html$Html$i,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$class('fa fa-fw fa-spin fa-circle-o-notch')
				]),
			_elm_lang$core$Native_List.fromArray(
				[])),
			_elm_lang$html$Html$text('Loading...')
		]));
var _concourse$atc$Login$routeWithRedirect = F2(
	function (redirect, route) {
		return A2(
			_elm_lang$core$Basics_ops['++'],
			route,
			function () {
				var _p5 = redirect;
				if (_p5 === '') {
					return '';
				} else {
					return A2(_elm_lang$core$Basics_ops['++'], '?redirect=', redirect);
				}
			}());
	});
var _concourse$atc$Login$teamSelectionRoute = function (redirect) {
	return A2(_concourse$atc$Login$routeWithRedirect, redirect, '/login');
};
var _concourse$atc$Login$loginRoute = F2(
	function (redirect, teamName) {
		return A2(
			_concourse$atc$Login$routeWithRedirect,
			redirect,
			A2(
				_elm_lang$core$Basics_ops['++'],
				'teams/',
				A2(_elm_lang$core$Basics_ops['++'], teamName, '/login')));
	});
var _concourse$atc$Login$PageWithRedirect = F2(
	function (a, b) {
		return {page: a, redirect: b};
	});
var _concourse$atc$Login$TeamSelectionModel = F3(
	function (a, b, c) {
		return {teamFilter: a, teams: b, redirect: c};
	});
var _concourse$atc$Login$LoginModel = F4(
	function (a, b, c, d) {
		return {teamName: a, authMethods: b, hasTeamSelectionInBrowserHistory: c, redirect: d};
	});
var _concourse$atc$Login$LoginPage = function (a) {
	return {ctor: 'LoginPage', _0: a};
};
var _concourse$atc$Login$TeamSelectionPage = {ctor: 'TeamSelectionPage'};
var _concourse$atc$Login$defaultPage = {page: _concourse$atc$Login$TeamSelectionPage, redirect: ''};
var _concourse$atc$Login$Login = function (a) {
	return {ctor: 'Login', _0: a};
};
var _concourse$atc$Login$TeamSelection = function (a) {
	return {ctor: 'TeamSelection', _0: a};
};
var _concourse$atc$Login$update = F2(
	function (action, model) {
		var _p6 = action;
		switch (_p6.ctor) {
			case 'Noop':
				return {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none};
			case 'FilterTeams':
				var _p7 = model;
				if (_p7.ctor === 'TeamSelection') {
					return {
						ctor: '_Tuple2',
						_0: _concourse$atc$Login$TeamSelection(
							_elm_lang$core$Native_Utils.update(
								_p7._0,
								{teamFilter: _p6._0})),
						_1: _elm_lang$core$Platform_Cmd$none
					};
				} else {
					return {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none};
				}
			case 'TeamsFetched':
				if (_p6._0.ctor === 'Ok') {
					var _p8 = model;
					if (_p8.ctor === 'TeamSelection') {
						return {
							ctor: '_Tuple2',
							_0: _concourse$atc$Login$TeamSelection(
								_elm_lang$core$Native_Utils.update(
									_p8._0,
									{
										teams: _elm_lang$core$Maybe$Just(_p6._0._0)
									})),
							_1: _elm_lang$core$Platform_Cmd$none
						};
					} else {
						return {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none};
					}
				} else {
					return A2(
						_elm_lang$core$Debug$log,
						A2(
							_elm_lang$core$Basics_ops['++'],
							'failed to fetch teams: ',
							_elm_lang$core$Basics$toString(_p6._0._0)),
						{ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none});
				}
			case 'SelectTeam':
				var _p9 = model;
				if (_p9.ctor === 'TeamSelection') {
					return {
						ctor: '_Tuple2',
						_0: model,
						_1: _elm_lang$navigation$Navigation$newUrl(
							A2(_concourse$atc$Login$loginRoute, _p9._0.redirect, _p6._0))
					};
				} else {
					return {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none};
				}
			case 'AuthFetched':
				if (_p6._0.ctor === 'Ok') {
					var _p10 = model;
					if (_p10.ctor === 'Login') {
						return {
							ctor: '_Tuple2',
							_0: _concourse$atc$Login$Login(
								_elm_lang$core$Native_Utils.update(
									_p10._0,
									{
										authMethods: _elm_lang$core$Maybe$Just(_p6._0._0)
									})),
							_1: _elm_lang$core$Platform_Cmd$none
						};
					} else {
						return {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none};
					}
				} else {
					return A2(
						_elm_lang$core$Debug$log,
						A2(
							_elm_lang$core$Basics_ops['++'],
							'failed to fetch auth methods: ',
							_elm_lang$core$Basics$toString(_p6._0._0)),
						{ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none});
				}
			default:
				var _p11 = model;
				if (_p11.ctor === 'Login') {
					var _p13 = _p11._0;
					var _p12 = _p13.hasTeamSelectionInBrowserHistory;
					if (_p12 === true) {
						return {
							ctor: '_Tuple2',
							_0: model,
							_1: _elm_lang$navigation$Navigation$back(1)
						};
					} else {
						return {
							ctor: '_Tuple2',
							_0: model,
							_1: _elm_lang$navigation$Navigation$newUrl(
								_concourse$atc$Login$teamSelectionRoute(_p13.redirect))
						};
					}
				} else {
					return {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none};
				}
		}
	});
var _concourse$atc$Login$GoBack = {ctor: 'GoBack'};
var _concourse$atc$Login$viewLogin = function (model) {
	return A2(
		_elm_lang$html$Html$div,
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html_Attributes$class('centered-contents')
			]),
		_elm_lang$core$Native_List.fromArray(
			[
				A2(
				_elm_lang$html$Html$div,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('small-title')
					]),
				_elm_lang$core$Native_List.fromArray(
					[
						A2(
						_elm_lang$html$Html$a,
						_elm_lang$core$Native_List.fromArray(
							[
								_concourse$atc$Login$onClickPreventDefault(_concourse$atc$Login$GoBack),
								_elm_lang$html$Html_Attributes$href(
								_concourse$atc$Login$teamSelectionRoute(model.redirect))
							]),
						_elm_lang$core$Native_List.fromArray(
							[
								A2(
								_elm_lang$html$Html$i,
								_elm_lang$core$Native_List.fromArray(
									[
										_elm_lang$html$Html_Attributes$class('fa fa-fw fa-chevron-left')
									]),
								_elm_lang$core$Native_List.fromArray(
									[])),
								_elm_lang$html$Html$text('back to team selection')
							]))
					])),
				A2(
				_elm_lang$html$Html$div,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('login-box auth-methods')
					]),
				A2(
					_elm_lang$core$Basics_ops['++'],
					_elm_lang$core$Native_List.fromArray(
						[
							A2(
							_elm_lang$html$Html$div,
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html_Attributes$class('centered-contents auth-methods-title')
								]),
							_elm_lang$core$Native_List.fromArray(
								[
									_elm_lang$html$Html$text('logging in to '),
									A2(
									_elm_lang$html$Html$span,
									_elm_lang$core$Native_List.fromArray(
										[
											_elm_lang$html$Html_Attributes$class('bright-text')
										]),
									_elm_lang$core$Native_List.fromArray(
										[
											_elm_lang$html$Html$text(model.teamName)
										]))
								]))
						]),
					function () {
						var _p14 = model.authMethods;
						if (_p14.ctor === 'Nothing') {
							return _elm_lang$core$Native_List.fromArray(
								[_concourse$atc$Login$viewLoading]);
						} else {
							var _p16 = _p14._0;
							var _p15 = {
								ctor: '_Tuple2',
								_0: _concourse$atc$Login$viewBasicAuthForm(_p16),
								_1: _concourse$atc$Login$viewOAuthButtons(_p16)
							};
							if (_p15._0.ctor === 'Just') {
								if (_p15._1.ctor === 'Just') {
									return _elm_lang$core$Native_List.fromArray(
										[_p15._1._0, _concourse$atc$Login$viewOrBar, _p15._0._0]);
								} else {
									return _elm_lang$core$Native_List.fromArray(
										[_p15._0._0]);
								}
							} else {
								if (_p15._1.ctor === 'Just') {
									return _elm_lang$core$Native_List.fromArray(
										[_p15._1._0]);
								} else {
									return _elm_lang$core$Native_List.fromArray(
										[_concourse$atc$Login$viewNoAuthButton]);
								}
							}
						}
					}()))
			]));
};
var _concourse$atc$Login$AuthFetched = function (a) {
	return {ctor: 'AuthFetched', _0: a};
};
var _concourse$atc$Login$SelectTeam = function (a) {
	return {ctor: 'SelectTeam', _0: a};
};
var _concourse$atc$Login$viewTeam = F2(
	function (redirect, team) {
		return A2(
			_elm_lang$html$Html$a,
			_elm_lang$core$Native_List.fromArray(
				[
					_concourse$atc$Login$onClickPreventDefault(
					_concourse$atc$Login$SelectTeam(team.name)),
					_elm_lang$html$Html_Attributes$href(
					A2(_concourse$atc$Login$loginRoute, redirect, team.name))
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html$text(team.name)
				]));
	});
var _concourse$atc$Login$TeamsFetched = function (a) {
	return {ctor: 'TeamsFetched', _0: a};
};
var _concourse$atc$Login$init = function (pageResult) {
	var pageWithRedirect = A2(_elm_lang$core$Result$withDefault, _concourse$atc$Login$defaultPage, pageResult);
	var _p17 = pageWithRedirect.page;
	if (_p17.ctor === 'TeamSelectionPage') {
		return {
			ctor: '_Tuple2',
			_0: _concourse$atc$Login$TeamSelection(
				{teamFilter: '', teams: _elm_lang$core$Maybe$Nothing, redirect: pageWithRedirect.redirect}),
			_1: A2(
				_elm_lang$core$Platform_Cmd$map,
				_concourse$atc$Login$TeamsFetched,
				A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _concourse$atc$Concourse_Team$fetchTeams))
		};
	} else {
		var _p18 = _p17._0;
		return {
			ctor: '_Tuple2',
			_0: _concourse$atc$Login$Login(
				{teamName: _p18, authMethods: _elm_lang$core$Maybe$Nothing, hasTeamSelectionInBrowserHistory: false, redirect: pageWithRedirect.redirect}),
			_1: A2(
				_elm_lang$core$Platform_Cmd$map,
				_concourse$atc$Login$AuthFetched,
				A3(
					_elm_lang$core$Task$perform,
					_elm_lang$core$Result$Err,
					_elm_lang$core$Result$Ok,
					_concourse$atc$Concourse_AuthMethod$fetchAuthMethods(_p18)))
		};
	}
};
var _concourse$atc$Login$urlUpdate = F2(
	function (pageResult, model) {
		var pageWithRedirect = A2(_elm_lang$core$Result$withDefault, _concourse$atc$Login$defaultPage, pageResult);
		var _p19 = pageWithRedirect.page;
		if (_p19.ctor === 'TeamSelectionPage') {
			return {
				ctor: '_Tuple2',
				_0: _concourse$atc$Login$TeamSelection(
					{teamFilter: '', teams: _elm_lang$core$Maybe$Nothing, redirect: pageWithRedirect.redirect}),
				_1: A2(
					_elm_lang$core$Platform_Cmd$map,
					_concourse$atc$Login$TeamsFetched,
					A3(_elm_lang$core$Task$perform, _elm_lang$core$Result$Err, _elm_lang$core$Result$Ok, _concourse$atc$Concourse_Team$fetchTeams))
			};
		} else {
			var _p20 = _p19._0;
			return {
				ctor: '_Tuple2',
				_0: _concourse$atc$Login$Login(
					{teamName: _p20, authMethods: _elm_lang$core$Maybe$Nothing, hasTeamSelectionInBrowserHistory: true, redirect: pageWithRedirect.redirect}),
				_1: A2(
					_elm_lang$core$Platform_Cmd$map,
					_concourse$atc$Login$AuthFetched,
					A3(
						_elm_lang$core$Task$perform,
						_elm_lang$core$Result$Err,
						_elm_lang$core$Result$Ok,
						_concourse$atc$Concourse_AuthMethod$fetchAuthMethods(_p20)))
			};
		}
	});
var _concourse$atc$Login$FilterTeams = function (a) {
	return {ctor: 'FilterTeams', _0: a};
};
var _concourse$atc$Login$Noop = {ctor: 'Noop'};
var _concourse$atc$Login$viewTeamSelection = function (model) {
	var filteredTeams = A2(
		_concourse$atc$Login$filterTeams,
		model.teamFilter,
		A2(
			_elm_lang$core$Maybe$withDefault,
			_elm_lang$core$Native_List.fromArray(
				[]),
			model.teams));
	return A2(
		_elm_lang$html$Html$div,
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html_Attributes$class('centered-contents')
			]),
		_elm_lang$core$Native_List.fromArray(
			[
				A2(
				_elm_lang$html$Html$div,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('small-title')
					]),
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html$text('select a team to login')
					])),
				A2(
				_elm_lang$html$Html$div,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('login-box team-selection')
					]),
				_elm_lang$core$Native_List.fromArray(
					[
						A2(
						_elm_lang$html$Html$form,
						_elm_lang$core$Native_List.fromArray(
							[
								_elm_lang$html$Html_Events$onSubmit(
								function () {
									var _p21 = {
										ctor: '_Tuple2',
										_0: _elm_lang$core$List$head(filteredTeams),
										_1: model.teamFilter
									};
									if (_p21._0.ctor === 'Nothing') {
										return _concourse$atc$Login$Noop;
									} else {
										if (_p21._1 === '') {
											return _concourse$atc$Login$Noop;
										} else {
											return _concourse$atc$Login$SelectTeam(_p21._0._0.name);
										}
									}
								}()),
								_elm_lang$html$Html_Attributes$class('filter-form input-holder')
							]),
						_elm_lang$core$Native_List.fromArray(
							[
								A2(
								_elm_lang$html$Html$i,
								_elm_lang$core$Native_List.fromArray(
									[
										_elm_lang$html$Html_Attributes$class('fa fa-fw fa-search')
									]),
								_elm_lang$core$Native_List.fromArray(
									[])),
								A2(
								_elm_lang$html$Html$input,
								_elm_lang$core$Native_List.fromArray(
									[
										_elm_lang$html$Html_Attributes$placeholder('filter teams'),
										_elm_lang$html$Html_Attributes$autofocus(true),
										_elm_lang$html$Html_Events$onInput(_concourse$atc$Login$FilterTeams)
									]),
								_elm_lang$core$Native_List.fromArray(
									[]))
							])),
						function () {
						var _p22 = model.teams;
						if (_p22.ctor === 'Nothing') {
							return _concourse$atc$Login$viewLoading;
						} else {
							return A2(
								_elm_lang$html$Html$div,
								_elm_lang$core$Native_List.fromArray(
									[]),
								A2(
									_elm_lang$core$List$map,
									_concourse$atc$Login$viewTeam(model.redirect),
									filteredTeams));
						}
					}()
					]))
			]));
};
var _concourse$atc$Login$view = function (model) {
	var _p23 = model;
	if (_p23.ctor === 'TeamSelection') {
		return _concourse$atc$Login$viewTeamSelection(_p23._0);
	} else {
		return _concourse$atc$Login$viewLogin(_p23._0);
	}
};

var _evancz$url_parser$UrlParser$oneOfHelp = F3(
	function (choices, chunks, formatter) {
		oneOfHelp:
		while (true) {
			var _p0 = choices;
			if (_p0.ctor === '[]') {
				return _elm_lang$core$Result$Err('Tried many parsers, but none of them worked!');
			} else {
				var _p1 = A2(_p0._0._0, chunks, formatter);
				if (_p1.ctor === 'Err') {
					var _v2 = _p0._1,
						_v3 = chunks,
						_v4 = formatter;
					choices = _v2;
					chunks = _v3;
					formatter = _v4;
					continue oneOfHelp;
				} else {
					return _elm_lang$core$Result$Ok(_p1._0);
				}
			}
		}
	});
var _evancz$url_parser$UrlParser$Chunks = F2(
	function (a, b) {
		return {seen: a, rest: b};
	});
var _evancz$url_parser$UrlParser$parse = F3(
	function (input, _p2, url) {
		var _p3 = _p2;
		var _p4 = A2(
			_p3._0,
			A2(
				_evancz$url_parser$UrlParser$Chunks,
				_elm_lang$core$Native_List.fromArray(
					[]),
				A2(_elm_lang$core$String$split, '/', url)),
			input);
		if (_p4.ctor === 'Err') {
			return _elm_lang$core$Result$Err(_p4._0);
		} else {
			var _p7 = _p4._0._1;
			var _p6 = _p4._0._0.rest;
			var _p5 = _p6;
			if (_p5.ctor === '[]') {
				return _elm_lang$core$Result$Ok(_p7);
			} else {
				if ((_p5._0 === '') && (_p5._1.ctor === '[]')) {
					return _elm_lang$core$Result$Ok(_p7);
				} else {
					return _elm_lang$core$Result$Err(
						A2(
							_elm_lang$core$Basics_ops['++'],
							'The parser worked, but /',
							A2(
								_elm_lang$core$Basics_ops['++'],
								A2(_elm_lang$core$String$join, '/', _p6),
								' was left over.')));
				}
			}
		}
	});
var _evancz$url_parser$UrlParser$Parser = function (a) {
	return {ctor: 'Parser', _0: a};
};
var _evancz$url_parser$UrlParser$s = function (str) {
	return _evancz$url_parser$UrlParser$Parser(
		F2(
			function (_p8, result) {
				var _p9 = _p8;
				var _p12 = _p9.rest;
				var _p10 = _p12;
				if (_p10.ctor === '[]') {
					return _elm_lang$core$Result$Err(
						A2(_elm_lang$core$Basics_ops['++'], 'Got to the end of the URL but wanted /', str));
				} else {
					var _p11 = _p10._0;
					return _elm_lang$core$Native_Utils.eq(_p11, str) ? _elm_lang$core$Result$Ok(
						{
							ctor: '_Tuple2',
							_0: A2(
								_evancz$url_parser$UrlParser$Chunks,
								A2(_elm_lang$core$List_ops['::'], _p11, _p9.seen),
								_p10._1),
							_1: result
						}) : _elm_lang$core$Result$Err(
						A2(
							_elm_lang$core$Basics_ops['++'],
							'Wanted /',
							A2(
								_elm_lang$core$Basics_ops['++'],
								str,
								A2(
									_elm_lang$core$Basics_ops['++'],
									' but got /',
									A2(_elm_lang$core$String$join, '/', _p12)))));
				}
			}));
};
var _evancz$url_parser$UrlParser$custom = F2(
	function (tipe, stringToSomething) {
		return _evancz$url_parser$UrlParser$Parser(
			F2(
				function (_p13, func) {
					var _p14 = _p13;
					var _p15 = _p14.rest;
					if (_p15.ctor === '[]') {
						return _elm_lang$core$Result$Err(
							A2(_elm_lang$core$Basics_ops['++'], 'Got to the end of the URL but wanted /', tipe));
					} else {
						var _p17 = _p15._0;
						var _p16 = stringToSomething(_p17);
						if (_p16.ctor === 'Ok') {
							return _elm_lang$core$Result$Ok(
								{
									ctor: '_Tuple2',
									_0: A2(
										_evancz$url_parser$UrlParser$Chunks,
										A2(_elm_lang$core$List_ops['::'], _p17, _p14.seen),
										_p15._1),
									_1: func(_p16._0)
								});
						} else {
							return _elm_lang$core$Result$Err(
								A2(
									_elm_lang$core$Basics_ops['++'],
									'Parsing `',
									A2(
										_elm_lang$core$Basics_ops['++'],
										_p17,
										A2(_elm_lang$core$Basics_ops['++'], '` went wrong: ', _p16._0))));
						}
					}
				}));
	});
var _evancz$url_parser$UrlParser$string = A2(_evancz$url_parser$UrlParser$custom, 'STRING', _elm_lang$core$Result$Ok);
var _evancz$url_parser$UrlParser$int = A2(_evancz$url_parser$UrlParser$custom, 'NUMBER', _elm_lang$core$String$toInt);
var _evancz$url_parser$UrlParser_ops = _evancz$url_parser$UrlParser_ops || {};
_evancz$url_parser$UrlParser_ops['</>'] = F2(
	function (_p19, _p18) {
		var _p20 = _p19;
		var _p21 = _p18;
		return _evancz$url_parser$UrlParser$Parser(
			F2(
				function (chunks, func) {
					return A2(
						_elm_lang$core$Result$andThen,
						A2(_p20._0, chunks, func),
						function (_p22) {
							var _p23 = _p22;
							return A2(_p21._0, _p23._0, _p23._1);
						});
				}));
	});
var _evancz$url_parser$UrlParser$oneOf = function (choices) {
	return _evancz$url_parser$UrlParser$Parser(
		_evancz$url_parser$UrlParser$oneOfHelp(choices));
};
var _evancz$url_parser$UrlParser$format = F2(
	function (input, _p24) {
		var _p25 = _p24;
		return _evancz$url_parser$UrlParser$Parser(
			F2(
				function (chunks, func) {
					var _p26 = A2(_p25._0, chunks, input);
					if (_p26.ctor === 'Err') {
						return _elm_lang$core$Result$Err(_p26._0);
					} else {
						return _elm_lang$core$Result$Ok(
							{
								ctor: '_Tuple2',
								_0: _p26._0._0,
								_1: func(_p26._0._1)
							});
					}
				}));
	});

var _concourse$atc$LoginPage$pageParser = _evancz$url_parser$UrlParser$oneOf(
	_elm_lang$core$Native_List.fromArray(
		[
			A2(
			_evancz$url_parser$UrlParser$format,
			_concourse$atc$Login$TeamSelectionPage,
			_evancz$url_parser$UrlParser$s('login')),
			A2(
			_evancz$url_parser$UrlParser$format,
			_concourse$atc$Login$LoginPage,
			A2(
				_evancz$url_parser$UrlParser_ops['</>'],
				_evancz$url_parser$UrlParser$s('teams'),
				A2(
					_evancz$url_parser$UrlParser_ops['</>'],
					_evancz$url_parser$UrlParser$string,
					_evancz$url_parser$UrlParser$s('login'))))
		]));
var _concourse$atc$LoginPage$redirectInserter = F2(
	function (uri, page) {
		return {page: page, redirect: uri};
	});
var _concourse$atc$LoginPage$pathnameParser = function (location) {
	var redirect = function () {
		var _p0 = A3(
			_elm_lang$core$Regex$find,
			_elm_lang$core$Regex$AtMost(1),
			_elm_lang$core$Regex$regex('[&?]redirect=([^&]*)'),
			location.search);
		if ((_p0.ctor === '::') && (_p0._1.ctor === '[]')) {
			return A2(
				_elm_lang$core$Maybe$withDefault,
				'',
				A2(
					_elm_lang$core$Maybe$withDefault,
					_elm_lang$core$Maybe$Nothing,
					_elm_lang$core$List$head(_p0._0.submatches)));
		} else {
			return '';
		}
	}();
	return A3(
		_evancz$url_parser$UrlParser$parse,
		_concourse$atc$LoginPage$redirectInserter(redirect),
		_concourse$atc$LoginPage$pageParser,
		A2(_elm_lang$core$String$dropLeft, 1, location.pathname));
};
var _concourse$atc$LoginPage$main = {
	main: A2(
		_elm_lang$navigation$Navigation$program,
		_elm_lang$navigation$Navigation$makeParser(_concourse$atc$LoginPage$pathnameParser),
		{
			init: _concourse$atc$Login$init,
			update: _concourse$atc$Login$update,
			urlUpdate: _concourse$atc$Login$urlUpdate,
			view: _concourse$atc$Login$view,
			subscriptions: _elm_lang$core$Basics$always(_elm_lang$core$Platform_Sub$none)
		})
};

var _elm_lang$animation_frame$Native_AnimationFrame = function()
{

var hasStartTime =
	window.performance &&
	window.performance.timing &&
	window.performance.timing.navigationStart;

var navStart = hasStartTime
	? window.performance.timing.navigationStart
	: Date.now();

var rAF = _elm_lang$core$Native_Scheduler.nativeBinding(function(callback)
{
	var id = requestAnimationFrame(function(time) {
		var timeNow = time
			? (time > navStart ? time : time + navStart)
			: Date.now();

		callback(_elm_lang$core$Native_Scheduler.succeed(timeNow));
	});

	return function() {
		cancelAnimationFrame(id);
	};
});

return {
	rAF: rAF
};

}();

var _elm_lang$animation_frame$AnimationFrame$rAF = _elm_lang$animation_frame$Native_AnimationFrame.rAF;
var _elm_lang$animation_frame$AnimationFrame$subscription = _elm_lang$core$Native_Platform.leaf('AnimationFrame');
var _elm_lang$animation_frame$AnimationFrame$State = F3(
	function (a, b, c) {
		return {subs: a, request: b, oldTime: c};
	});
var _elm_lang$animation_frame$AnimationFrame$init = _elm_lang$core$Task$succeed(
	A3(
		_elm_lang$animation_frame$AnimationFrame$State,
		_elm_lang$core$Native_List.fromArray(
			[]),
		_elm_lang$core$Maybe$Nothing,
		0));
var _elm_lang$animation_frame$AnimationFrame$onEffects = F3(
	function (router, subs, _p0) {
		var _p1 = _p0;
		var _p5 = _p1.request;
		var _p4 = _p1.oldTime;
		var _p2 = {ctor: '_Tuple2', _0: _p5, _1: subs};
		if (_p2._0.ctor === 'Nothing') {
			if (_p2._1.ctor === '[]') {
				return _elm_lang$core$Task$succeed(
					A3(
						_elm_lang$animation_frame$AnimationFrame$State,
						_elm_lang$core$Native_List.fromArray(
							[]),
						_elm_lang$core$Maybe$Nothing,
						_p4));
			} else {
				return A2(
					_elm_lang$core$Task$andThen,
					_elm_lang$core$Process$spawn(
						A2(
							_elm_lang$core$Task$andThen,
							_elm_lang$animation_frame$AnimationFrame$rAF,
							_elm_lang$core$Platform$sendToSelf(router))),
					function (pid) {
						return A2(
							_elm_lang$core$Task$andThen,
							_elm_lang$core$Time$now,
							function (time) {
								return _elm_lang$core$Task$succeed(
									A3(
										_elm_lang$animation_frame$AnimationFrame$State,
										subs,
										_elm_lang$core$Maybe$Just(pid),
										time));
							});
					});
			}
		} else {
			if (_p2._1.ctor === '[]') {
				return A2(
					_elm_lang$core$Task$andThen,
					_elm_lang$core$Process$kill(_p2._0._0),
					function (_p3) {
						return _elm_lang$core$Task$succeed(
							A3(
								_elm_lang$animation_frame$AnimationFrame$State,
								_elm_lang$core$Native_List.fromArray(
									[]),
								_elm_lang$core$Maybe$Nothing,
								_p4));
					});
			} else {
				return _elm_lang$core$Task$succeed(
					A3(_elm_lang$animation_frame$AnimationFrame$State, subs, _p5, _p4));
			}
		}
	});
var _elm_lang$animation_frame$AnimationFrame$onSelfMsg = F3(
	function (router, newTime, _p6) {
		var _p7 = _p6;
		var _p10 = _p7.subs;
		var diff = newTime - _p7.oldTime;
		var send = function (sub) {
			var _p8 = sub;
			if (_p8.ctor === 'Time') {
				return A2(
					_elm_lang$core$Platform$sendToApp,
					router,
					_p8._0(newTime));
			} else {
				return A2(
					_elm_lang$core$Platform$sendToApp,
					router,
					_p8._0(diff));
			}
		};
		return A2(
			_elm_lang$core$Task$andThen,
			_elm_lang$core$Process$spawn(
				A2(
					_elm_lang$core$Task$andThen,
					_elm_lang$animation_frame$AnimationFrame$rAF,
					_elm_lang$core$Platform$sendToSelf(router))),
			function (pid) {
				return A2(
					_elm_lang$core$Task$andThen,
					_elm_lang$core$Task$sequence(
						A2(_elm_lang$core$List$map, send, _p10)),
					function (_p9) {
						return _elm_lang$core$Task$succeed(
							A3(
								_elm_lang$animation_frame$AnimationFrame$State,
								_p10,
								_elm_lang$core$Maybe$Just(pid),
								newTime));
					});
			});
	});
var _elm_lang$animation_frame$AnimationFrame$Diff = function (a) {
	return {ctor: 'Diff', _0: a};
};
var _elm_lang$animation_frame$AnimationFrame$diffs = function (tagger) {
	return _elm_lang$animation_frame$AnimationFrame$subscription(
		_elm_lang$animation_frame$AnimationFrame$Diff(tagger));
};
var _elm_lang$animation_frame$AnimationFrame$Time = function (a) {
	return {ctor: 'Time', _0: a};
};
var _elm_lang$animation_frame$AnimationFrame$times = function (tagger) {
	return _elm_lang$animation_frame$AnimationFrame$subscription(
		_elm_lang$animation_frame$AnimationFrame$Time(tagger));
};
var _elm_lang$animation_frame$AnimationFrame$subMap = F2(
	function (func, sub) {
		var _p11 = sub;
		if (_p11.ctor === 'Time') {
			return _elm_lang$animation_frame$AnimationFrame$Time(
				function (_p12) {
					return func(
						_p11._0(_p12));
				});
		} else {
			return _elm_lang$animation_frame$AnimationFrame$Diff(
				function (_p13) {
					return func(
						_p11._0(_p13));
				});
		}
	});
_elm_lang$core$Native_Platform.effectManagers['AnimationFrame'] = {pkg: 'elm-lang/animation-frame', init: _elm_lang$animation_frame$AnimationFrame$init, onEffects: _elm_lang$animation_frame$AnimationFrame$onEffects, onSelfMsg: _elm_lang$animation_frame$AnimationFrame$onSelfMsg, tag: 'sub', subMap: _elm_lang$animation_frame$AnimationFrame$subMap};

var _concourse$atc$Pipeline$outputNodes = F2(
	function (job, _p0) {
		var _p1 = _p0;
		return _elm_lang$core$Native_List.fromArray(
			[]);
	});
var _concourse$atc$Pipeline$viewOutputNode = function (resourceName) {
	return A2(
		_elm_lang$html$Html$a,
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html_Attributes$href('#')
			]),
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html$text(resourceName)
			]));
};
var _concourse$atc$Pipeline$viewConstrainedInputNode = function (resourceName) {
	return A2(
		_elm_lang$html$Html$a,
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html_Attributes$href('#')
			]),
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html$text(resourceName)
			]));
};
var _concourse$atc$Pipeline$viewInputNode = function (resourceName) {
	return A2(
		_elm_lang$html$Html$a,
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html_Attributes$href('#')
			]),
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html$text(resourceName)
			]));
};
var _concourse$atc$Pipeline$jobResources = function (_p2) {
	var _p3 = _p2;
	return _elm_lang$core$Set$size(
		_elm_lang$core$Set$fromList(
			A2(
				_elm_lang$core$Basics_ops['++'],
				A2(
					_elm_lang$core$List$map,
					function (_) {
						return _.resource;
					},
					_p3.inputs),
				A2(
					_elm_lang$core$List$map,
					function (_) {
						return _.resource;
					},
					_p3.outputs))));
};
var _concourse$atc$Pipeline$viewJobNode = function (job) {
	var linkAttrs = function () {
		var _p4 = {ctor: '_Tuple2', _0: job.finishedBuild, _1: job.nextBuild};
		if (_p4._0.ctor === 'Just') {
			if (_p4._1.ctor === 'Just') {
				return _elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class(
						A2(
							_elm_lang$core$Basics_ops['++'],
							_concourse$atc$Concourse_BuildStatus$show(_p4._0._0.status),
							' started')),
						_elm_lang$html$Html_Attributes$href(_p4._1._0.url)
					]);
			} else {
				var _p5 = _p4._0._0;
				return _elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class(
						_concourse$atc$Concourse_BuildStatus$show(_p5.status)),
						_elm_lang$html$Html_Attributes$href(_p5.url)
					]);
			}
		} else {
			if (_p4._1.ctor === 'Just') {
				return _elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('no-builds started'),
						_elm_lang$html$Html_Attributes$href(_p4._1._0.url)
					]);
			} else {
				return _elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('no-builds'),
						_elm_lang$html$Html_Attributes$href(job.url)
					]);
			}
		}
	}();
	return A2(
		_elm_lang$html$Html$a,
		linkAttrs,
		_elm_lang$core$Native_List.fromArray(
			[
				_elm_lang$html$Html$text(job.name)
			]));
};
var _concourse$atc$Pipeline$viewNode = function (_p6) {
	var _p7 = _p6;
	var idAttr = _elm_lang$html$Html_Attributes$id(
		A2(
			_elm_lang$core$Basics_ops['++'],
			'node-',
			_elm_lang$core$Basics$toString(_p7.id)));
	var _p8 = _p7.label;
	switch (_p8.ctor) {
		case 'JobNode':
			return A2(
				_elm_lang$html$Html$div,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('node job'),
						idAttr
					]),
				_elm_lang$core$Native_List.fromArray(
					[
						_concourse$atc$Pipeline$viewJobNode(_p8._0)
					]));
		case 'InputNode':
			return A2(
				_elm_lang$html$Html$div,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('node input'),
						idAttr
					]),
				_elm_lang$core$Native_List.fromArray(
					[
						_concourse$atc$Pipeline$viewInputNode(_p8._0.resourceName)
					]));
		case 'ConstrainedInputNode':
			return A2(
				_elm_lang$html$Html$div,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('node input constrained'),
						idAttr
					]),
				_elm_lang$core$Native_List.fromArray(
					[
						_concourse$atc$Pipeline$viewConstrainedInputNode(_p8._0.resourceName)
					]));
		default:
			return A2(
				_elm_lang$html$Html$div,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('node output'),
						idAttr
					]),
				_elm_lang$core$Native_List.fromArray(
					[
						_concourse$atc$Pipeline$viewOutputNode(_p8._0.resourceName)
					]));
	}
};
var _concourse$atc$Pipeline$viewGrid = function (grid) {
	var _p9 = grid;
	switch (_p9.ctor) {
		case 'Cell':
			return _concourse$atc$Pipeline$viewNode(_p9._0.node);
		case 'Serial':
			return A2(
				_elm_lang$html$Html$div,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('serial-grid')
					]),
				A2(
					_elm_lang$core$Basics_ops['++'],
					_concourse$atc$Pipeline$viewSerial(_p9._0),
					_concourse$atc$Pipeline$viewSerial(_p9._1)));
		case 'Parallel':
			return A2(
				_elm_lang$html$Html$div,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('parallel-grid')
					]),
				A2(_elm_lang$core$List$map, _concourse$atc$Pipeline$viewGrid, _p9._0));
		default:
			return _elm_lang$html$Html$text('');
	}
};
var _concourse$atc$Pipeline$viewSerial = function (grid) {
	var _p10 = grid;
	if (_p10.ctor === 'Serial') {
		return A2(
			_elm_lang$core$Basics_ops['++'],
			_concourse$atc$Pipeline$viewSerial(_p10._0),
			_concourse$atc$Pipeline$viewSerial(_p10._1));
	} else {
		return _elm_lang$core$Native_List.fromArray(
			[
				_concourse$atc$Pipeline$viewGrid(grid)
			]);
	}
};
var _concourse$atc$Pipeline$nodeHeight = function (_p11) {
	var _p12 = _p11;
	var _p13 = _p12.label;
	if (_p13.ctor === 'JobNode') {
		return A2(
			_elm_lang$core$Basics$max,
			1,
			_concourse$atc$Pipeline$jobResources(_p13._0));
	} else {
		return 1;
	}
};
var _concourse$atc$Pipeline$viewMatrixCell = function (mnode) {
	var _p14 = mnode;
	switch (_p14.ctor) {
		case 'MatrixSpacer':
			return A2(
				_elm_lang$html$Html$td,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$class('spacer')
					]),
				_elm_lang$core$Native_List.fromArray(
					[]));
		case 'MatrixNode':
			var _p15 = _p14._0.node;
			return A2(
				_elm_lang$html$Html$td,
				_elm_lang$core$Native_List.fromArray(
					[
						_elm_lang$html$Html_Attributes$rowspan(
						_concourse$atc$Pipeline$nodeHeight(_p15))
					]),
				_elm_lang$core$Native_List.fromArray(
					[
						_concourse$atc$Pipeline$viewNode(_p15)
					]));
		default:
			return _elm_lang$html$Html$text('');
	}
};
var _concourse$atc$Pipeline$viewRow = function (row) {
	return A2(
		_elm_lang$html$Html$tr,
		_elm_lang$core$Native_List.fromArray(
			[]),
		A2(_elm_lang$core$List$map, _concourse$atc$Pipeline$viewMatrixCell, row));
};
var _concourse$atc$Pipeline$view = function (model) {
	var _p16 = model.error;
	if (_p16.ctor === 'Just') {
		return _elm_lang$html$Html$text(
			A2(_elm_lang$core$Basics_ops['++'], 'error: ', _p16._0));
	} else {
		return A2(
			_elm_lang$html$Html$div,
			_elm_lang$core$Native_List.fromArray(
				[
					_elm_lang$html$Html_Attributes$class('pipeline-grid')
				]),
			_elm_lang$core$Native_List.fromArray(
				[
					_concourse$atc$Pipeline$viewGrid(
					_concourse$atc$Grid$fromGraph(model.graph))
				]));
	}
};
var _concourse$atc$Pipeline$Model = F6(
	function (a, b, c, d, e, f) {
		return {fit: a, pipelineLocator: b, groups: c, jobs: d, graph: e, error: f};
	});
var _concourse$atc$Pipeline$Flags = F3(
	function (a, b, c) {
		return {teamName: a, pipelineName: b, groups: c};
	});
var _concourse$atc$Pipeline$ConstrainedInputNode = function (a) {
	return {ctor: 'ConstrainedInputNode', _0: a};
};
var _concourse$atc$Pipeline$constrainedInputNode = F4(
	function (jobs, resourceName, dependentJob, upstreamJobName) {
		return _concourse$atc$Pipeline$ConstrainedInputNode(
			{
				resourceName: resourceName,
				dependentJob: dependentJob,
				upstreamJob: A2(_elm_lang$core$Dict$get, upstreamJobName, jobs)
			});
	});
var _concourse$atc$Pipeline$OutputNode = function (a) {
	return {ctor: 'OutputNode', _0: a};
};
var _concourse$atc$Pipeline$InputNode = function (a) {
	return {ctor: 'InputNode', _0: a};
};
var _concourse$atc$Pipeline$inputNodes = F3(
	function (jobs, job, _p17) {
		var _p18 = _p17;
		var _p20 = _p18.resource;
		var _p19 = _p18.passed;
		return _elm_lang$core$List$isEmpty(_p19) ? _elm_lang$core$Native_List.fromArray(
			[
				_concourse$atc$Pipeline$InputNode(
				{resourceName: _p20, dependentJob: job})
			]) : A2(
			_elm_lang$core$List$map,
			A3(_concourse$atc$Pipeline$constrainedInputNode, jobs, _p20, job),
			_p19);
	});
var _concourse$atc$Pipeline$jobResourceNodes = F2(
	function (jobs, job) {
		return A2(
			_elm_lang$core$Basics_ops['++'],
			A2(
				_elm_lang$core$List$concatMap,
				A2(_concourse$atc$Pipeline$inputNodes, jobs, job),
				job.inputs),
			A2(
				_elm_lang$core$List$concatMap,
				_concourse$atc$Pipeline$outputNodes(job),
				job.outputs));
	});
var _concourse$atc$Pipeline$JobNode = function (a) {
	return {ctor: 'JobNode', _0: a};
};
var _concourse$atc$Pipeline$jobId = F2(
	function (nodes, job) {
		var _p22 = A2(
			_elm_lang$core$List$filter,
			function (_p21) {
				return A2(
					F2(
						function (x, y) {
							return _elm_lang$core$Native_Utils.eq(x, y);
						}),
					_concourse$atc$Pipeline$JobNode(job),
					function (_) {
						return _.label;
					}(_p21));
			},
			nodes);
		if (_p22.ctor === '::') {
			return _p22._0.id;
		} else {
			return _elm_lang$core$Native_Utils.crashCase(
				'Pipeline',
				{
					start: {line: 323, column: 3},
					end: {line: 328, column: 52}
				},
				_p22)('impossible: job index not found');
		}
	});
var _concourse$atc$Pipeline$nodeEdges = F2(
	function (allNodes, _p24) {
		var _p25 = _p24;
		var _p28 = _p25.id;
		var _p26 = _p25.label;
		switch (_p26.ctor) {
			case 'JobNode':
				return _elm_lang$core$Native_List.fromArray(
					[]);
			case 'InputNode':
				return _elm_lang$core$Native_List.fromArray(
					[
						A3(
						_elm_community$graph$Graph$Edge,
						_p28,
						A2(_concourse$atc$Pipeline$jobId, allNodes, _p26._0.dependentJob),
						{ctor: '_Tuple0'})
					]);
			case 'ConstrainedInputNode':
				return A2(
					_elm_lang$core$List_ops['::'],
					A3(
						_elm_community$graph$Graph$Edge,
						_p28,
						A2(_concourse$atc$Pipeline$jobId, allNodes, _p26._0.dependentJob),
						{ctor: '_Tuple0'}),
					function () {
						var _p27 = _p26._0.upstreamJob;
						if (_p27.ctor === 'Just') {
							return _elm_lang$core$Native_List.fromArray(
								[
									A3(
									_elm_community$graph$Graph$Edge,
									A2(_concourse$atc$Pipeline$jobId, allNodes, _p27._0),
									_p28,
									{ctor: '_Tuple0'})
								]);
						} else {
							return _elm_lang$core$Native_List.fromArray(
								[]);
						}
					}());
			default:
				return _elm_lang$core$Native_List.fromArray(
					[
						A3(
						_elm_community$graph$Graph$Edge,
						A2(_concourse$atc$Pipeline$jobId, allNodes, _p26._0.upstreamJob),
						_p28,
						{ctor: '_Tuple0'})
					]);
		}
	});
var _concourse$atc$Pipeline$initGraph = function (jobs) {
	var jobsByName = A3(
		_elm_lang$core$List$foldl,
		F2(
			function (job, dict) {
				return A3(_elm_lang$core$Dict$insert, job.name, job, dict);
			}),
		_elm_lang$core$Dict$empty,
		jobs);
	var resourceNodes = A2(
		_elm_lang$core$List$concatMap,
		_concourse$atc$Pipeline$jobResourceNodes(jobsByName),
		jobs);
	var jobNodes = A2(_elm_lang$core$List$map, _concourse$atc$Pipeline$JobNode, jobs);
	var graphNodes = A2(
		_elm_lang$core$List$indexedMap,
		_elm_community$graph$Graph$Node,
		_elm_lang$core$List$concat(
			_elm_lang$core$Native_List.fromArray(
				[jobNodes, resourceNodes])));
	return A2(
		_elm_community$graph$Graph$fromNodesAndEdges,
		graphNodes,
		A2(
			_elm_lang$core$List$concatMap,
			_concourse$atc$Pipeline$nodeEdges(graphNodes),
			graphNodes));
};
var _concourse$atc$Pipeline$update = F2(
	function (msg, model) {
		var _p29 = msg;
		switch (_p29.ctor) {
			case 'Noop':
				return {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none};
			case 'JobsFetched':
				if (_p29._0.ctor === 'Ok') {
					var _p31 = _p29._0._0;
					var filtered = _elm_lang$core$Set$isEmpty(model.groups) ? _p31 : A2(
						_elm_lang$core$List$filter,
						function (_p30) {
							return A2(
								_elm_lang$core$List$any,
								A2(_elm_lang$core$Basics$flip, _elm_lang$core$Set$member, model.groups),
								function (_) {
									return _.groups;
								}(_p30));
						},
						_p31);
					return {
						ctor: '_Tuple2',
						_0: _elm_lang$core$Native_Utils.update(
							model,
							{
								jobs: filtered,
								graph: _concourse$atc$Pipeline$initGraph(filtered)
							}),
						_1: _elm_lang$core$Platform_Cmd$none
					};
				} else {
					return {
						ctor: '_Tuple2',
						_0: _elm_lang$core$Native_Utils.update(
							model,
							{
								error: _elm_lang$core$Maybe$Just(
									_elm_lang$core$Basics$toString(_p29._0._0))
							}),
						_1: _elm_lang$core$Platform_Cmd$none
					};
				}
			default:
				var _p32 = model.fit;
				if (_p32.ctor === 'Just') {
					return {
						ctor: '_Tuple2',
						_0: _elm_lang$core$Native_Utils.update(
							model,
							{fit: _elm_lang$core$Maybe$Nothing}),
						_1: _p32._0(
							{ctor: '_Tuple0'})
					};
				} else {
					return {ctor: '_Tuple2', _0: model, _1: _elm_lang$core$Platform_Cmd$none};
				}
		}
	});
var _concourse$atc$Pipeline$Frame = {ctor: 'Frame'};
var _concourse$atc$Pipeline$subscriptions = function (model) {
	return (_elm_lang$core$List$isEmpty(model.jobs) || _elm_lang$core$Native_Utils.eq(model.fit, _elm_lang$core$Maybe$Nothing)) ? _elm_lang$core$Platform_Sub$none : _elm_lang$animation_frame$AnimationFrame$times(
		_elm_lang$core$Basics$always(_concourse$atc$Pipeline$Frame));
};
var _concourse$atc$Pipeline$JobsFetched = function (a) {
	return {ctor: 'JobsFetched', _0: a};
};
var _concourse$atc$Pipeline$fetchJobs = function (locator) {
	return A2(
		_elm_lang$core$Platform_Cmd$map,
		_concourse$atc$Pipeline$JobsFetched,
		A3(
			_elm_lang$core$Task$perform,
			_elm_lang$core$Result$Err,
			_elm_lang$core$Result$Ok,
			_concourse$atc$Concourse_Job$fetchJobs(locator)));
};
var _concourse$atc$Pipeline$init = F2(
	function (fit, flags) {
		var model = {
			fit: _elm_lang$core$Maybe$Just(fit),
			pipelineLocator: {teamName: flags.teamName, pipelineName: flags.pipelineName},
			groups: _elm_lang$core$Set$fromList(flags.groups),
			jobs: _elm_lang$core$Native_List.fromArray(
				[]),
			graph: _elm_community$graph$Graph$empty,
			error: _elm_lang$core$Maybe$Nothing
		};
		return {
			ctor: '_Tuple2',
			_0: model,
			_1: _concourse$atc$Pipeline$fetchJobs(model.pipelineLocator)
		};
	});
var _concourse$atc$Pipeline$Noop = {ctor: 'Noop'};

var _concourse$atc$PipelinePage$fit = _elm_lang$core$Native_Platform.outgoingPort(
	'fit',
	function (v) {
		return null;
	});
var _concourse$atc$PipelinePage$main = {
	main: _elm_lang$html$Html_App$programWithFlags(
		{
			init: _concourse$atc$Pipeline$init(_concourse$atc$PipelinePage$fit),
			update: _concourse$atc$Pipeline$update,
			view: _concourse$atc$Pipeline$view,
			subscriptions: _concourse$atc$Pipeline$subscriptions
		}),
	flags: A2(
		_elm_lang$core$Json_Decode$andThen,
		A2(
			_elm_lang$core$Json_Decode_ops[':='],
			'groups',
			_elm_lang$core$Json_Decode$list(_elm_lang$core$Json_Decode$string)),
		function (groups) {
			return A2(
				_elm_lang$core$Json_Decode$andThen,
				A2(_elm_lang$core$Json_Decode_ops[':='], 'pipelineName', _elm_lang$core$Json_Decode$string),
				function (pipelineName) {
					return A2(
						_elm_lang$core$Json_Decode$andThen,
						A2(_elm_lang$core$Json_Decode_ops[':='], 'teamName', _elm_lang$core$Json_Decode$string),
						function (teamName) {
							return _elm_lang$core$Json_Decode$succeed(
								{groups: groups, pipelineName: pipelineName, teamName: teamName});
						});
				});
		})
};

var Elm = {};
Elm['BuildPage'] = Elm['BuildPage'] || {};
_elm_lang$core$Native_Platform.addPublicModule(Elm['BuildPage'], 'BuildPage', typeof _concourse$atc$BuildPage$main === 'undefined' ? null : _concourse$atc$BuildPage$main);
Elm['JobPage'] = Elm['JobPage'] || {};
_elm_lang$core$Native_Platform.addPublicModule(Elm['JobPage'], 'JobPage', typeof _concourse$atc$JobPage$main === 'undefined' ? null : _concourse$atc$JobPage$main);
Elm['LoginPage'] = Elm['LoginPage'] || {};
_elm_lang$core$Native_Platform.addPublicModule(Elm['LoginPage'], 'LoginPage', typeof _concourse$atc$LoginPage$main === 'undefined' ? null : _concourse$atc$LoginPage$main);
Elm['PipelinePage'] = Elm['PipelinePage'] || {};
_elm_lang$core$Native_Platform.addPublicModule(Elm['PipelinePage'], 'PipelinePage', typeof _concourse$atc$PipelinePage$main === 'undefined' ? null : _concourse$atc$PipelinePage$main);

if (typeof define === "function" && define['amd'])
{
  define([], function() { return Elm; });
  return;
}

if (typeof module === "object")
{
  module['exports'] = Elm;
  return;
}

var globalElm = this['Elm'];
if (typeof globalElm === "undefined")
{
  this['Elm'] = Elm;
  return;
}

for (var publicModule in Elm)
{
  if (publicModule in globalElm)
  {
    throw new Error('There are two Elm modules called `' + publicModule + '` on this page! Rename one of them.');
  }
  globalElm[publicModule] = Elm[publicModule];
}

}).call(this);

