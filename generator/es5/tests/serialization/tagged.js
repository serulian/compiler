$module('tagged', function () {
  var $static = this;
  this.$struct('bad44f92', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        somefield: SomeField,
      };
      instance.$markruntimecreated();
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'SomeField', 'somefield', function () {
      return $g.____testlib.basictypes.Integer;
    }, function () {
      return $g.____testlib.basictypes.Integer;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|29dc432d<bad44f92>": true,
        "equals|4|29dc432d<43834c3f>": true,
        "Stringify|2|29dc432d<5cffd9b5>": true,
        "Mapping|2|29dc432d<df58fcbd<any>>": true,
        "Clone|2|29dc432d<bad44f92>": true,
        "String|2|29dc432d<5cffd9b5>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = $t.markpromising(function () {
    var $result;
    var jsonString;
    var s;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            s = $g.tagged.SomeStruct.new($t.fastbox(2, $g.____testlib.basictypes.Integer));
            jsonString = $t.fastbox('{"somefield":2}', $g.____testlib.basictypes.String);
            $promise.maybe(s.Stringify($g.____testlib.basictypes.JSON)()).then(function ($result0) {
              $result = $g.____testlib.basictypes.String.$equals($result0, jsonString);
              $current = 1;
              $continue($resolve, $reject);
              return;
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
  });
});
