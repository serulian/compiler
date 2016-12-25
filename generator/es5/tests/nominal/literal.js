$module('literal', function () {
  var $static = this;
  $static.TEST = function () {
    return $t.box('foo' == 'foo', $g.____testlib.basictypes.Boolean);
  };
});
