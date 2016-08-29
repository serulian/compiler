$module('attributes', function () {
  var $static = this;
  $static.SimpleFunction = function (props) {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            props.$index($t.box("data-foo", $g.____testlib.basictypes.String)).then(function ($result2) {
              return $g.____testlib.basictypes.String.$equals($t.cast($result2, $g.____testlib.basictypes.String, false), $t.box("bar", $g.____testlib.basictypes.String)).then(function ($result1) {
                return $promise.resolve($t.unbox($result1)).then(function ($result0) {
                  return ($promise.shortcircuit($result0, true) || props.$index($t.box("an-attr-here", $g.____testlib.basictypes.String))).then(function ($result3) {
                    $result = $t.box($result0 && $t.unbox($t.cast($result3, $g.____testlib.basictypes.Boolean, false)), $g.____testlib.basictypes.Boolean);
                    $current = 1;
                    $continue($resolve, $reject);
                    return;
                  });
                });
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.Mapping($t.any).overObject({
              "data-foo": $t.box("bar", $g.____testlib.basictypes.String),
              "an-attr-here": $t.box(true, $g.____testlib.basictypes.Boolean),
            }).then(function ($result1) {
              return $g.attributes.SimpleFunction($result1).then(function ($result0) {
                $result = $result0;
                $current = 1;
                $continue($resolve, $reject);
                return;
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
