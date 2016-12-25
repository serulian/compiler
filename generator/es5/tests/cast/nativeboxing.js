$module('nativeboxing', function () {
  var $static = this;
  $static.TEST = function () {
    var r;
    var s;
    var sany;
    s = 'hello world';
    sany = s;
    r = $t.fastbox($t.cast(sany, $global.String, false), $g.____testlib.basictypes.String);
    return $g.____testlib.basictypes.String.$equals(r, $t.fastbox('hello world', $g.____testlib.basictypes.String));
  };
});
