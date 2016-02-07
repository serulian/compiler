$module('string', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $t.nominalwrap('hello world', $g.____testlib.basictypes.String);
            $t.nominalwrap("hi world", $g.____testlib.basictypes.String);
            $t.nominalwrap('single quote with "quoted"', $g.____testlib.basictypes.String);
            $t.nominalwrap("double quote with 'quoted'", $g.____testlib.basictypes.String);
            $t.nominalwrap("escaped \" quote", $g.____testlib.basictypes.String);
            $t.nominalwrap('escaped \' quote', $g.____testlib.basictypes.String);
            $state.current = -1;
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
