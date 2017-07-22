$module('literal', function () {
  var $static = this;
  $static.TEST = function () {
    return $t.box('foo' == 'foo', $g.________testlib.basictypes.Boolean);
  };
});
