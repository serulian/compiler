$module('classprops', function () {
  var $static = this;
  this.$class('c89a612d', 'SomeProps', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (BoolValue, StringValue) {
      var instance = new $static();
      instance.BoolValue = BoolValue;
      instance.StringValue = StringValue;
      return $promise.resolve(instance);
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.SimpleFunction = function (props) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.String.$equals(props.StringValue, $t.fastbox("hello world", $g.____testlib.basictypes.String)).then(function ($result2) {
              return $promise.resolve($result2.$wrapped).then(function ($result1) {
                return $promise.resolve($result1 && props.BoolValue.$wrapped).then(function ($result0) {
                  $result = $t.fastbox($result0 && !(props.OptionalValue == null), $g.____testlib.basictypes.Boolean);
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
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
    var $temp0;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.classprops.SomeProps.new($t.fastbox(true, $g.____testlib.basictypes.Boolean), $t.fastbox("hello world", $g.____testlib.basictypes.String)).then(function ($result1) {
              $temp0 = $result1;
              return $g.classprops.SimpleFunction(($temp0, $temp0.OptionalValue = $t.fastbox(42, $g.____testlib.basictypes.Integer), $temp0)).then(function ($result0) {
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
