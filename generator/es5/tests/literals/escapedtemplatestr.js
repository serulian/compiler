$module('escapedtemplatestr', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.____testlib.basictypes.Slice($g.____testlib.basictypes.String).overArray([$t.nominalwrap("hello 'world'! \"This is a\n\tlong quote!\"", $g.____testlib.basictypes.String)]).then(function ($result0) {
              return $g.____testlib.basictypes.Slice($g.____testlib.basictypes.Stringable).overArray([]).then(function ($result1) {
                return $g.____testlib.basictypes.formatTemplateString($result0, $result1).then(function ($result2) {
                  $result = $result2;
                  $state.current = 1;
                  $callback($state);
                });
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            $result;
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
