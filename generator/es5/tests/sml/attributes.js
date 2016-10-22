$module('attributes', function () {
  var $static = this;
  $static.SimpleFunction = function (props) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            props.$index($t.fastbox("data-foo", $g.____testlib.basictypes.String)).then(function ($result2) {
              return $g.____testlib.basictypes.String.$equals($t.cast($result2, $g.____testlib.basictypes.String, false), $t.fastbox("bar", $g.____testlib.basictypes.String)).then(function ($result1) {
                return $promise.resolve($result1.$wrapped).then(function ($result0) {
                  return ($promise.shortcircuit($result0, true) || props.$index($t.fastbox("an-attr-here", $g.____testlib.basictypes.String))).then(function ($result3) {
                    $result = $t.fastbox($result0 && $t.cast($result3, $g.____testlib.basictypes.Boolean, false).$wrapped, $g.____testlib.basictypes.Boolean);
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
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.Mapping($t.any).overObject({
              "an-attr-here": $t.fastbox(true, $g.____testlib.basictypes.Boolean),
              "data-foo": $t.fastbox("bar", $g.____testlib.basictypes.String),
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
