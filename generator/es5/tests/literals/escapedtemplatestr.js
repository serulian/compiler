$module('escapedtemplatestr', function () {
  var $static = this;
  $static.DoSomething = function () {
    $t.fastbox("hello 'world'! \"This is a\n\tlong quote!\"", $g.________testlib.basictypes.String);
    return;
  };
});
