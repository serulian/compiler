$module('mappingliteral', function () {
  var $static = this;
  $static.TEST = function () {
    return $g.________testlib.basictypes.Mapping($g.________testlib.basictypes.Boolean).overObject(function () {
      var obj = {
      };
      obj['somekey'] = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
      obj['anotherkey'] = $t.fastbox(false, $g.________testlib.basictypes.Boolean);
      return obj;
    }()).$index($t.fastbox('somekey', $g.________testlib.basictypes.String));
  };
});
